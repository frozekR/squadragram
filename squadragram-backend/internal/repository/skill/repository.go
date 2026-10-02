package skill

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
	return skills, nil
}
