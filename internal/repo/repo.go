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
func (r *repository) GetBalanceByUserID(ctx context.Context, userID int64) (models.Balance, error) {
	balanceDB := models.BalanceDB{}

	query := `
		SELECT
			b.id,
			b.amount,
			b.deleted_at,
			b.created_at
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
	err := row.Scan(
		&balanceDB.ID,
		&balanceDB.Amount,
		&balanceDB.CreatedAt,
		&balanceDB.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Balance{}, nil
		}
		return models.Balance{}, err
	}

	return balanceDB.ToModelBalance(), nil
}

func (r *repository) GetBalanceAmountByUserID(ctx context.Context, userID int64) (int64, error) {
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
		return nil, fmt.Errorf("failed to query GetBalanceHistoryByUserID: %w", err)
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
			return nil, fmt.Errorf("failed to scan GetBalanceHistoryByUserID: %w", err)
		}
		balanceHistoryDB = append(balanceHistoryDB, bh)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read rows GetBalanceHistoryByUserID: %w", err)
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
		return nil, fmt.Errorf("failed to query GetInventoryMerchItems: %w", err)
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
			return nil, fmt.Errorf("failed to scan GetInventoryMerchItems: %w", err)
		}
		inventoryMerchDB = append(inventoryMerchDB, im)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read rows GetInventoryMerchItems: %w", err)
	}

	inventoryMerch := make([]models.InventoryMerch, 0, len(inventoryMerchDB))
	for _, e := range inventoryMerchDB {
		inventoryMerch = append(inventoryMerch, e.ToModelInventoryMerch())
	}

	return inventoryMerch, nil
}

func (r *repository) BuyItemTX(ctx context.Context, userID, balanceID, inventoryID, merchID, price int64, username, item string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// списание баланса
	query := `
		UPDATE
			shop."balance"
		SET
			amount = amount - $1
		WHERE
			id = $2
	`
	cmdTag, err := tx.Exec(ctx, query, price, balanceID)
	if err != nil {
		return fmt.Errorf("failed to execute query BuyItemTX: %v", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no balance rows updated BuyItemTX")
	}

	// создание записи в истории транзакций
	query = `
		INSERT INTO
			shop."balance_history" (balance_id, transaction_amount, sender, recipient)
		VALUES
			($1, $2, $3, 'AvitoShop')
	`
	cmdTag, err = tx.Exec(ctx, query, balanceID, price, username)
	if err != nil {
		return fmt.Errorf("failed to execute query BuyItemTX: %v", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows inserted balance history BuyItemTX")
	}

	// создание записи-связки для инвентаря с данным предметом
	query = `
		INSERT INTO
			shop."inventory_merch" (inventory_id, merch_id, name, count)
		VALUES
			($1, $2, $3, 1)
		ON CONFLICT (inventory_id, merch_id)
		DO UPDATE SET count = shop."inventory_merch".count + 1
	`
	cmdTag, err = tx.Exec(ctx, query, inventoryID, merchID, item)
	if err != nil {
		return fmt.Errorf("failed to execute query BuyItemTX: %v", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows inserted inventory merch BuyItemTX")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction BuyItemTX: %w", err)
	}

	return nil
}

func (r *repository) GetInventoryIDByUserID(ctx context.Context, userID int64) (int64, error) {
	var inventoryID int64

	query := `
		SELECT
			i.id
		FROM
			shop."inventory" i
		WHERE
			i.user_id = $1
	`

	row := r.db.QueryRow(ctx, query, userID)
	err := row.Scan(&inventoryID)
	if err != nil {
		return 0, err
	}

	return inventoryID, nil
}

// Merch
func (r *repository) GetMerchByName(ctx context.Context, name string) (models.Merch, error) {
	merchDB := models.MerchDB{}

	query := `
		SELECT
			m.id,
			m.name,
			m.price,
			m.deleted_at,
			m.created_at
		FROM
			shop."merch" m
		WHERE
			m.name = $1
	`

	row := r.db.QueryRow(ctx, query, name)
	err := row.Scan(
		&merchDB.ID,
		&merchDB.Name,
		&merchDB.Price,
		&merchDB.CreatedAt,
		&merchDB.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Merch{}, nil
		}
		return models.Merch{}, err
	}

	return merchDB.ToModelMerch(), nil
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
// 			return nil, fmt.Errorf("failed to scan merch data: %w", err)
// 		}
// 		merchesDB = append(merchesDB, m)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to read rows: %w", err)
// 	}

// 	merches := make([]models.Merch, 0, len(merchesDB))
// 	for _, e := range merchesDB {
// 		merches = append(merches, e.ToModelMerch())
// 	}

// 	return merches, nil
// }
