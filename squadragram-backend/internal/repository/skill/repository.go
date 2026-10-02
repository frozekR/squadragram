package skill

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

func (r *Repository) GetSkillByUUID(ctx context.Context, u uuid.UUID) (model.Skill, error) {
	query, args := QuerySelectByUUID(u)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.Skill{}, fmt.Errorf("query skill by UUID: %w", err)
	}
	skill, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Skill])
	if err != nil {
		return model.Skill{}, fmt.Errorf("collect skill by UUID: %w", err)
	}
	skill.MediaMetadata, err = r.media.GetMediaMetadataByOwner(ctx, model.OwnerTypeSkill, skill.UUID)
	if err != nil {
		return model.Skill{}, err
	}
	return skill, nil
}

func (r *Repository) GetSkills(ctx context.Context) ([]model.Skill, error) {
	query, args := QuerySelectAll()
	return r.collectSkills(ctx, query, args)
}

// GetSkillsByCharacterUUID возвращает скиллы, привязанные к персонажу через char_id.
// Для персонажа без скиллов или неизвестного UUID возвращается пустой список.
func (r *Repository) GetSkillsByCharacterUUID(ctx context.Context, characterUUID uuid.UUID) ([]model.Skill, error) {
	query, args := QuerySelectByCharacterUUID(characterUUID)
	return r.collectSkills(ctx, query, args)
}

// CreateSkill создаёт скилл только для существующего персонажа.
// Если CharacterUUID не найден, ошибка оборачивает pgx.ErrNoRows.
func (r *Repository) CreateSkill(ctx context.Context, dto model.CreateSkillDTO) (model.Skill, error) {
	query, args := QueryInsert(uuid.New(), dto)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.Skill{}, fmt.Errorf("insert skill: %w", err)
	}
	skill, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Skill])
	if err != nil {
		return model.Skill{}, fmt.Errorf("collect created skill for character %s: %w", dto.CharacterUUID, err)
	}
	skill.MediaMetadata = []model.MediaMetadata{}
	return skill, nil
}

func (r *Repository) collectSkills(ctx context.Context, query string, args []any) ([]model.Skill, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query skills: %w", err)
	}
	skills, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Skill])
	if err != nil {
		return nil, fmt.Errorf("collect skills: %w", err)
	}
	// Связь с медиа идёт по UUID скилла, отдельно от его связи char_id с персонажем.
	owners := make([]uuid.UUID, len(skills))
	for i := range skills {
		owners[i] = skills[i].UUID
	}
	media, err := r.media.GetByOwners(ctx, model.OwnerTypeSkill, owners)
	if err != nil {
		return nil, err
	}
	for i := range skills {
		skills[i].MediaMetadata = media[skills[i].UUID]
		if skills[i].MediaMetadata == nil {
			skills[i].MediaMetadata = []model.MediaMetadata{}
		}
	}
	return skills, nil
}

func (r *Repository) UpdateSkill(ctx context.Context, u uuid.UUID, dto model.CreateSkillDTO) (model.Skill, error) {
	query, args := QueryUpdate(u, dto)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return model.Skill{}, fmt.Errorf("update skill: %w", err)
	}
	skill, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Skill])
	if err != nil {
		return model.Skill{}, fmt.Errorf("collect updated skill: %w", err)
	}
	skill.MediaMetadata, err = r.media.GetMediaMetadataByOwner(ctx, model.OwnerTypeSkill, skill.UUID)
	return skill, err
}
