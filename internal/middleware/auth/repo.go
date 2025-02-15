package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/devWaylander/coins_store/pkg/log"
	"github.com/devWaylander/coins_store/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) txRollback(ctx context.Context, tx pgx.Tx, err error) {
	if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
		err = fmt.Errorf("failed to rollback transaction: %w", err)
		log.Logger.Err(rollbackErr).Msg(err.Error())
	}
}

func (r *repository) CreateUserTX(ctx context.Context, username, passwordHash string) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	// создание баланса
	var balanceID int64
	query := `INSERT INTO shop."balance" (amount) VALUES (1000) RETURNING id`
	err = tx.QueryRow(ctx, query).Scan(&balanceID)
	if err != nil {
		r.txRollback(ctx, tx, err)
		return 0, fmt.Errorf("failed to create balance CreateUserTX: %w", err)
	}

	// создание пользователя
	var userID int64
	query = `
		INSERT INTO 
			shop."user" (balance_id, username, password_hash)
		VALUES 
			($1, $2, $3)
		RETURNING 
			id
	`
	err = tx.QueryRow(ctx, query, balanceID, username, passwordHash).Scan(&userID)
	if err != nil {
		r.txRollback(ctx, tx, err)
		return 0, fmt.Errorf("failed to create user CreateUserTX: %w", err)
	}

	// создание инвентаря пользователя
	var inventoryID int64
	query = `INSERT INTO shop."inventory" (user_id) VALUES ($1) RETURNING id`
	err = tx.QueryRow(ctx, query, userID).Scan(&inventoryID)
	if err != nil {
		r.txRollback(ctx, tx, err)
		return 0, fmt.Errorf("failed to create inventory CreateUserTX: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		err = fmt.Errorf("failed to commit transaction CreateUserTX: %w", err)
		r.txRollback(ctx, tx, err)
		return 0, err
	}

	return userID, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := models.UserDB{}

	query := `
		SELECT
			u.id,
			u.balance_id,
			u.username,
			u.password_hash,
			u.created_at,
			u.deleted_at
		FROM
			shop."user" u
		WHERE
			u.username = $1 AND u.deleted_at IS NULL
	`

	row := r.db.QueryRow(ctx, query, username)
	err := row.Scan(
		&user.ID,
		&user.BalanceID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.User{}, nil
		}
		return &models.User{}, fmt.Errorf("GetUserByUsername failed: %w", err)
	}

	return user.ToModelUser(), nil
}

func (r *repository) GetUserPassHashByUsername(ctx context.Context, username string) (string, error) {
	passHash := ""

	query := `
		SELECT
			u.password_hash
		FROM
			shop."user" u
		WHERE
			u.username = $1 AND u.deleted_at IS NULL
	`

	row := r.db.QueryRow(ctx, query, username)
	err := row.Scan(&passHash)
	if err != nil {
		return "", fmt.Errorf("GetUserPassHashByUsername failed: %w", err)
	}

	return passHash, nil
}
