package character

import (
	"context"
	"squadraton-backend/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCharacterByUUIDQuery(ctx context.Context, uuid uuid.UUID) (model.Character, error) {
	query, args := QuerySelectByUUID(uuid)
	rows, err := r.db.Query(ctx, query, args...)
	character, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Character])
	if err != nil {
		return model.Character{}, err
	}
	return character, nil
}
