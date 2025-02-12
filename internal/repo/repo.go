package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/devWaylander/coins_store/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

// Balance
func (r *repository) GetBalanceByUserID(ctx context.Context, userID int64) (int64, error) {
	amount := 0

	query := `
		SELECT
			b.amount
		FROM
			shop."user" u
		INNER JOIN
			shop."balance" b
		ON
			u.balance_id = b.id
		WHERE
			u.id = $1 AND u.deleted_at IS NULL AND b.deleted_at IS NULL
	`

	row := r.db.QueryRow(ctx, query, userID)
	err := row.Scan(&amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return int64(amount), nil
}

// Balance History
func (r *repository) GetBalanceHistoryByUserID(ctx context.Context, userID int64) ([]models.BalanceHistory, error) {
	var balanceHistoryDB []models.BalanceHistoryDB

	query := `
		SELECT
			bh.id,
			bh.balance_id,
			bh.transaction_amount,
			bh.sender,
			bh.recipient,
			bh.deleted_at,
			bh.created_at
		FROM
			shop."balance_history" bh
		INNER JOIN
			shop."user" u
		ON
			u.balance_id = bh.balance_id
		WHERE
			u.id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		// TODO: вынести ошибки
		return nil, fmt.Errorf("failed to query balance history data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		bh := models.BalanceHistoryDB{}
		if err := rows.Scan(
			&bh.ID,
			&bh.BalanceID,
			&bh.TransactionAmount,
			&bh.Sender,
			&bh.Recipient,
			&bh.DeletedAt,
			&bh.CreatedAt,
		); err != nil {
			// TODO: вынести ошибки
			return nil, fmt.Errorf("failed to scan balance history data: %w", err)
		}
		balanceHistoryDB = append(balanceHistoryDB, bh)
	}

	if err := rows.Err(); err != nil {
		// TODO: вынести ошибки
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	balanceHistory := make([]models.BalanceHistory, 0, len(balanceHistoryDB))
	for _, e := range balanceHistoryDB {
		balanceHistory = append(balanceHistory, e.ToModelBalanceHistory())
	}

	return balanceHistory, nil
}

// Inventory
func (r *repository) GetInventoryMerchItems(ctx context.Context, userID int64) ([]models.InventoryMerch, error) {
	var inventoryMerchDB []models.InventoryMerchDB

	query := `
		SELECT 
			im.inventory_id, 
			im.merch_id, 
			im.name,
			im.count, 
			im.deleted_at, 
			im.created_at
		FROM 
			shop."inventory_merch" im
		INNER JOIN 
			shop."inventory" i 
		ON 
			im.inventory_id = i.id
		WHERE 
			i.user_id = $1 AND im.deleted_at IS NULL
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		// TODO: вынести ошибки
		return nil, fmt.Errorf("failed to query inventory merch data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		im := models.InventoryMerchDB{}
		if err := rows.Scan(
			&im.InventoryID,
			&im.MerchID,
			&im.Name,
			&im.Count,
			&im.DeletedAt,
			&im.CreatedAt,
		); err != nil {
			// TODO: вынести ошибки
			return nil, fmt.Errorf("failed to scan inventory merch data: %w", err)
		}
		inventoryMerchDB = append(inventoryMerchDB, im)
	}

	if err := rows.Err(); err != nil {
		// TODO: вынести ошибки
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	inventoryMerch := make([]models.InventoryMerch, 0, len(inventoryMerchDB))
	for _, e := range inventoryMerchDB {
		inventoryMerch = append(inventoryMerch, e.ToModelInventoryMerch())
	}

	return inventoryMerch, nil
}

// func (r *repository) GetInventoryMerchesByIDs(ctx context.Context, merchesIDs []int64) ([]models.Merch, error) {
// 	merchesDB := make([]models.MerchDB, 0, len(merchesIDs))

// 	query := `
// 		SELECT
// 			m.id,
// 			m.name,
// 			m.price,
// 			m.deleted_at,
// 			m.created_at
// 		FROM
// 			shop."merch" m
// 		WHERE
// 			m.id = ANY($1)
// 	`

// 	rows, err := r.db.Query(ctx, query, merchesIDs)
// 	if err != nil {
// 		// TODO: вынести ошибки
// 		return nil, fmt.Errorf("failed to query merch data: %w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		m := models.MerchDB{}
// 		if err := rows.Scan(
// 			&m.ID,
// 			&m.Name,
// 			&m.Price,
// 			&m.DeletedAt,
// 			&m.CreatedAt,
// 		); err != nil {
// 			// TODO: вынести ошибки
// 			return nil, fmt.Errorf("failed to scan merch data: %w", err)
// 		}
// 		merchesDB = append(merchesDB, m)
// 	}

// 	if err := rows.Err(); err != nil {
// 		// TODO: вынести ошибки
// 		return nil, fmt.Errorf("failed to read rows: %w", err)
// 	}

// 	merches := make([]models.Merch, 0, len(merchesDB))
// 	for _, e := range merchesDB {
// 		merches = append(merches, e.ToModelMerch())
// 	}

// 	return merches, nil
// }
