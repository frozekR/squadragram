package skill

import (
	"squadraton-backend/internal/model"

	"github.com/google/uuid"
)

func QuerySelectAll() (string, []any) {
	return `SELECT id, uuid, name, COALESCE(description, '') AS description, type, char_id, sort_order
	FROM skills
	ORDER BY sort_order, id`, []any{}
}

func QuerySelectByUUID(u uuid.UUID) (string, []any) {
	return `SELECT id, uuid, name, COALESCE(description, '') AS description, type, char_id, sort_order
	FROM skills
	WHERE uuid = $1`, []any{u}
}

func QuerySelectByCharacterUUID(characterUUID uuid.UUID) (string, []any) {
	// JOIN связывает skills.char_id с characters.id, а фильтр использует
	// публичный UUID персонажа. Результат содержит все его скиллы.
	return `SELECT s.id, s.uuid, s.name, COALESCE(s.description, '') AS description, s.type, s.char_id, s.sort_order
	FROM skills AS s
	JOIN characters AS c ON c.id = s.char_id
	WHERE c.uuid = $1
	ORDER BY s.sort_order, s.id`, []any{characterUUID}
}

func QueryInsert(u uuid.UUID, dto model.CreateSkillDTO) (string, []any) {
	// INSERT ... SELECT переводит UUID персонажа во внутренний id в одном запросе.
	// Если персонаж не найден, вставки не будет и RETURNING не вернёт строку.
	return `INSERT INTO skills (uuid, name, description, type, char_id, sort_order)
	SELECT $1, $2, $3, $4, c.id, COALESCE($6::integer, 0)
	FROM characters AS c
	WHERE c.uuid = $5
	RETURNING id, uuid, name, COALESCE(description, '') AS description, type, char_id, sort_order`,
		[]any{u, dto.Name, dto.Description, dto.SkillType, dto.CharacterUUID, dto.SortOrder}
}

func QueryUpdate(u uuid.UUID, dto model.CreateSkillDTO) (string, []any) {
	return `UPDATE skills AS s SET name = $2, description = $3, type = $4, char_id = c.id, sort_order = COALESCE($6::integer, s.sort_order)
	FROM characters AS c WHERE s.uuid = $1 AND c.uuid = $5
	RETURNING s.id, s.uuid, s.name, COALESCE(s.description, '') AS description, s.type, s.char_id, s.sort_order`,
		[]any{u, dto.Name, dto.Description, dto.SkillType, dto.CharacterUUID, dto.SortOrder}
}
