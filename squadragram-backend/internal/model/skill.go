package model

import (
	"github.com/google/uuid"
)

type Skill struct {
	ID          int       `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SkillType   SkillType `json:"type" db:"type"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	// CharacterID хранит characters.id: один персонаж может иметь несколько скиллов.
	// Указатель позволяет читать существующие записи с NULL в skills.char_id.
	CharacterID   *int            `json:"character_id" db:"char_id"`
	MediaMetadata []MediaMetadata `json:"media_metadata" db:"-"`
}

type CreateSkillDTO struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	SkillType   SkillType `json:"type" validate:"required"`
	// Клиент передаёт публичный UUID персонажа; репозиторий находит его id
	// и сохраняет его в skills.char_id, который защищён внешним ключом.
	CharacterUUID uuid.UUID `json:"character_uuid" validate:"required"`
	// Порядок внутри категории type. При равенстве используется стабильный id.
	// nil при обновлении сохраняет прежний порядок для совместимости клиентов.
	SortOrder *int `json:"sort_order,omitempty"`
}
