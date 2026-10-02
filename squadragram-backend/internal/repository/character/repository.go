package character

import (
	"context"
	"fmt"
	"squadraton-backend/internal/model"
	"squadraton-backend/internal/repository"
	mediametadata "squadraton-backend/internal/repository/media_metadata"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db    repository.DB
	media *mediametadata.Repository
}

func NewRepository(db repository.DB) *Repository {
	return &Repository{db: db, media: mediametadata.NewRepository(db)}
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
	character.MediaMetadata, err = r.media.GetMediaMetadataByOwner(ctx, model.OwnerTypeCharacter, character.UUID)
	if err != nil {
		return model.Character{}, err
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
	// Медиа загружаются пачкой по UUID: число запросов не растёт со списком персонажей.
	owners := make([]uuid.UUID, len(characters))
	for i := range characters {
		owners[i] = characters[i].UUID
	}
	media, err := r.media.GetByOwners(ctx, model.OwnerTypeCharacter, owners)
	if err != nil {
		return nil, err
	}
	for i := range characters {
		characters[i].MediaMetadata = media[characters[i].UUID]
		if characters[i].MediaMetadata == nil {
			characters[i].MediaMetadata = []model.MediaMetadata{}
		}
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
	character.MediaMetadata = []model.MediaMetadata{}
	return character, nil
}

func (r *Repository) UpdateCharacter(ctx context.Context, u uuid.UUID, dto model.CreateCharacterDTO) (model.Character, error) {
	query, args := QueryUpdate(u, dto)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.Character{}, fmt.Errorf("update character: %w", err)
	}
	character, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Character])
	if err != nil {
		return model.Character{}, fmt.Errorf("collect updated character: %w", err)
	}
	character.MediaMetadata, err = r.media.GetMediaMetadataByOwner(ctx, model.OwnerTypeCharacter, character.UUID)
	return character, err
}
