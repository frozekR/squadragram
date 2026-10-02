package mediametadata

import (
	"squadraton-backend/internal/model"

	"github.com/google/uuid"
)

func QueryInsert(metadata model.MediaMetadata) (string, []any) {
	return `INSERT INTO media_metadata (owner_type, owner_uuid, media_type, object_key)
	VALUES ($1, $2, $3, $4)
	RETURNING id, owner_type, owner_uuid, media_type, object_key`,
		[]any{metadata.OwnerType, metadata.OwnerUUID, metadata.MediaType, metadata.ObjectKey}
}

func QuerySelectByID(id int) (string, []any) {
	return `SELECT id, owner_type, owner_uuid, media_type, object_key
	FROM media_metadata WHERE id = $1`, []any{id}
}

func QuerySelectByOwners(ownerType model.OwnerType, owners []uuid.UUID) (string, []any) {
	// Один запрос загружает медиа всех UUID выбранного типа. Тип владельца
	// обязателен: UUID другой сущности не должен возвращать чужое медиа.
	return `SELECT id, owner_type, owner_uuid, media_type, object_key
	FROM media_metadata WHERE owner_type = $1 AND owner_uuid = ANY($2::uuid[])
	ORDER BY id`, []any{ownerType, owners}
}

func QueryReplace(id int, key string) (string, []any) {
	// Lock and capture the current key in the same statement as the replacement.
	// Concurrent replacements each clean up the object they actually replaced.
	return `WITH previous AS (
		SELECT id, object_key FROM media_metadata WHERE id = $1 FOR UPDATE
	)
	UPDATE media_metadata AS m SET object_key = $2
	FROM previous AS p WHERE m.id = p.id
	RETURNING m.id, m.owner_type, m.owner_uuid, m.media_type, m.object_key,
		p.object_key AS previous_object_key`, []any{id, key}
}

func QueryDelete(id int) (string, []any) {
	return `DELETE FROM media_metadata WHERE id = $1
	RETURNING id, owner_type, owner_uuid, media_type, object_key`, []any{id}
}
