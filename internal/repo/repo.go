package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) GetBalanceByUserID(ctx context.Context, userID int64) (int64, error) {
	return 0, nil
}
