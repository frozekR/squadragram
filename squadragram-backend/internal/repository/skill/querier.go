package skill

import (
	"squadraton-backend/internal/model"

	"github.com/google/uuid"
)

func QuerySelectAll() (string, []any) {
	return `SELECT id, uuid, name, COALESCE(description, '') AS description, type, char_id
	FROM skills
	ORDER BY id`, []any{}
}

func QuerySelectByUUID(u uuid.UUID) (string, []any) {
	return `SELECT id, uuid, name, COALESCE(description, '') AS description, type, char_id
	FROM skills
	WHERE uuid = $1`, []any{u}
}

func QuerySelectByCharacterUUID(characterUUID uuid.UUID) (string, []any) {
	// JOIN связывает skills.char_id с characters.id, а фильтр использует
	// публичный UUID персонажа. Результат содержит все его скиллы.
	return `SELECT s.id, s.uuid, s.name, COALESCE(s.description, '') AS description, s.type, s.char_id
	FROM skills AS s
	JOIN characters AS c ON c.id = s.char_id
	WHERE c.uuid = $1
	ORDER BY s.id`, []any{characterUUID}
}

func QueryInsert(u uuid.UUID, dto model.CreateSkillDTO) (string, []any) {
	// INSERT ... SELECT переводит UUID персонажа во внутренний id в одном запросе.
	// Если персонаж не найден, вставки не будет и RETURNING не вернёт строку.
	return `INSERT INTO skills (uuid, name, description, type, char_id)
	SELECT $1, $2, $3, $4, c.id
	FROM characters AS c
	WHERE c.uuid = $5
	RETURNING id, uuid, name, COALESCE(description, '') AS description, type, char_id`,
		[]any{u, dto.Name, dto.Description, dto.SkillType, dto.CharacterUUID}
}
