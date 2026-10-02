package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// DB is implemented by both pgxpool.Pool and pgx.Tx.
type DB interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}
