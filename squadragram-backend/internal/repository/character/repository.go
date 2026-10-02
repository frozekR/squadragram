package character

import (
	"context"
	"fmt"
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

func (r *Repository) GetCharacterByUUID(ctx context.Context, u uuid.UUID) (model.Character, error) {
	query, args := QuerySelectByUUID(u)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.Character{}, fmt.Errorf("query character by UUID: %w", err)
	}
	character, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Character])
	if err != nil {
		return model.Character{}, fmt.Errorf("collect character by UUID: %w", err)
	}
	return character, nil
}

func (r *Repository) GetCharacters(ctx context.Context) ([]model.Character, error) {
	query, args := QuerySelectAll()
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query characters: %w", err)
	}
	characters, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Character])
	if err != nil {
		return nil, fmt.Errorf("collect characters: %w", err)
	}
	return characters, nil
}

func (r *Repository) CreateCharacter(ctx context.Context, dto model.CreateCharacterDTO) (model.Character, error) {
	query, args := QueryInsert(uuid.New(), dto.Name, dto.Description, dto.Role)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.Character{}, fmt.Errorf("insert character: %w", err)
	}
	character, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Character])
	if err != nil {
		return model.Character{}, fmt.Errorf("collect created character: %w", err)
	}
	return character, nil
}
