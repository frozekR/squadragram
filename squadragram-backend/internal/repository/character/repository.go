package character

import (
	"context"

	"squadraton-backend/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, character model.Character) (model.Character, error) {
	if character.UUID == uuid.Nil {
		character.UUID = uuid.New()
	}

	err := r.db.QueryRow(ctx, `
		INSERT INTO characters (uuid, name, description, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, character.UUID, character.Name, character.Description, character.Role).Scan(&character.ID)
	if err != nil {
		return model.Character{}, err
	}

	return character, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]model.Character, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, uuid, name, description, role
		FROM characters
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []model.Character
	for rows.Next() {
		var character model.Character
		if err := rows.Scan(
			&character.ID,
			&character.UUID,
			&character.Name,
			&character.Description,
			&character.Role,
		); err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (model.Character, error) {
	var character model.Character
	err := r.db.QueryRow(ctx, `
		SELECT id, uuid, name, description, role
		FROM characters
		WHERE id = $1
	`, id).Scan(
		&character.ID,
		&character.UUID,
		&character.Name,
		&character.Description,
		&character.Role,
	)
	return character, err
}

func (r *Repository) GetByUUID(ctx context.Context, id uuid.UUID) (model.Character, error) {
	var character model.Character
	err := r.db.QueryRow(ctx, `
		SELECT id, uuid, name, description, role
		FROM characters
		WHERE uuid = $1
	`, id).Scan(
		&character.ID,
		&character.UUID,
		&character.Name,
		&character.Description,
		&character.Role,
	)
	return character, err
}

var _ Querier = (*Repository)(nil)
