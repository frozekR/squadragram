package character

import (
	"squadraton-backend/internal/model"

	"github.com/google/uuid"
)

func QuerySelectAll() (string, []any) {
	return `SELECT id, uuid, name, COALESCE(description, '') AS description, role
	FROM characters
	ORDER BY id`, []any{}
}

func QuerySelectByUUID(u uuid.UUID) (string, []any) {
	return `SELECT id, uuid, name, COALESCE(description, '') AS description, role
	 FROM characters
	 WHERE uuid = $1
	 `, []any{u}
}

func QueryInsert(u uuid.UUID, name string, desc string, role model.CharacterRole) (string, []any) {
	return `INSERT INTO characters (uuid, name, description, role)
	VALUES ($1, $2, $3, $4)
	RETURNING id, uuid, name, COALESCE(description, '') AS description, role
	`, []any{u, name, desc, role}
}

func QueryUpdate(u uuid.UUID, dto model.CreateCharacterDTO) (string, []any) {
	return `UPDATE characters SET name = $2, description = $3, role = $4
	WHERE uuid = $1
	RETURNING id, uuid, name, COALESCE(description, '') AS description, role`,
		[]any{u, dto.Name, dto.Description, dto.Role}
}
