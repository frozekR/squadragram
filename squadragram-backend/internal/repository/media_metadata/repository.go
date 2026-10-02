package mediametadata

import (
	"context"
	"errors"
	"fmt"
	"squadraton-backend/internal/model"
	"squadraton-backend/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrUnsupportedOwner = errors.New("unsupported owner type")

type mediaReplacement struct {
	model.MediaMetadata
	PreviousObjectKey string `db:"previous_object_key"`
}

type Repository struct {
	db          repository.DB
	ownerTables map[model.OwnerType]string
}

func NewRepository(db repository.DB) *Repository {
	return &Repository{db: db, ownerTables: map[model.OwnerType]string{
		model.OwnerTypeCharacter: "characters",
		model.OwnerTypeSkill:     "skills",
		model.OwnerTypeEmote:     "emotes",
	}}
}

// RegisterOwnerType вызывается при настройке приложения, до обработки запросов.
// Для новой сущности добавьте её таблицу с uuid и значение в PostgreSQL owner_type.
// Связь media_metadata не требует отдельного столбца для каждого владельца.
func (r *Repository) RegisterOwnerType(ownerType model.OwnerType, table string) {
	r.ownerTables[ownerType] = table
}

func (r *Repository) OwnerExists(ctx context.Context, ownerType model.OwnerType, u uuid.UUID) (bool, error) {
	table, ok := r.ownerTables[ownerType]
	if !ok {
		return false, ErrUnsupportedOwner
	}
	// Имена таблиц задаются кодом приложения; пользователь передаёт только тип и UUID.
	query := "SELECT EXISTS (SELECT 1 FROM " + pgx.Identifier{table}.Sanitize() + " WHERE uuid = $1)"
	var exists bool
	if err := r.db.QueryRow(ctx, query, u).Scan(&exists); err != nil {
		return false, fmt.Errorf("check media owner: %w", err)
	}
	return exists, nil
}

func (r *Repository) CreateMediaMetadata(ctx context.Context, metadata model.MediaMetadata) (model.MediaMetadata, error) {
	query, args := QueryInsert(metadata)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.MediaMetadata{}, fmt.Errorf("insert media metadata: %w", err)
	}
	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MediaMetadata])
	if err != nil {
		return model.MediaMetadata{}, fmt.Errorf("collect media metadata: %w", err)
	}
	return result, nil
}

func (r *Repository) GetMediaMetadataByID(ctx context.Context, id int) (model.MediaMetadata, error) {
	query, args := QuerySelectByID(id)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.MediaMetadata{}, fmt.Errorf("query media metadata: %w", err)
	}
	metadata, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MediaMetadata])
	if err != nil {
		return model.MediaMetadata{}, fmt.Errorf("collect media metadata: %w", err)
	}
	return metadata, nil
}

func (r *Repository) GetMediaMetadataByOwner(ctx context.Context, ownerType model.OwnerType, u uuid.UUID) ([]model.MediaMetadata, error) {
	byOwner, err := r.GetByOwners(ctx, ownerType, []uuid.UUID{u})
	if err != nil {
		return nil, err
	}
	if metadata := byOwner[u]; metadata != nil {
		return metadata, nil
	}
	return []model.MediaMetadata{}, nil
}

func (r *Repository) GetByOwners(ctx context.Context, ownerType model.OwnerType, owners []uuid.UUID) (map[uuid.UUID][]model.MediaMetadata, error) {
	result := make(map[uuid.UUID][]model.MediaMetadata)
	if len(owners) == 0 {
		return result, nil
	}
	query, args := QuerySelectByOwners(ownerType, owners)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query owners media metadata: %w", err)
	}
	metadata, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.MediaMetadata])
	if err != nil {
		return nil, fmt.Errorf("collect owners media metadata: %w", err)
	}
	for _, item := range metadata {
		result[item.OwnerUUID] = append(result[item.OwnerUUID], item)
	}
	return result, nil
}

func (r *Repository) ReplaceMediaMetadata(ctx context.Context, id int, key string) (model.MediaMetadata, string, error) {
	query, args := QueryReplace(id, key)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.MediaMetadata{}, "", fmt.Errorf("replace media: %w", err)
	}
	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[mediaReplacement])
	if err != nil {
		return model.MediaMetadata{}, "", fmt.Errorf("collect replaced media: %w", err)
	}
	return result.MediaMetadata, result.PreviousObjectKey, nil
}

func (r *Repository) DeleteMediaMetadata(ctx context.Context, id int) (model.MediaMetadata, error) {
	query, args := QueryDelete(id)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.MediaMetadata{}, fmt.Errorf("delete media: %w", err)
	}
	metadata, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.MediaMetadata])
	if err != nil {
		return model.MediaMetadata{}, fmt.Errorf("collect deleted media: %w", err)
	}
	return metadata, nil
}
