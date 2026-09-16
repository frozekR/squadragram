package model

import "github.com/google/uuid"

type Character struct {
	ID          int           `json:"id"`
	UUID        uuid.UUID     `json:"uuid"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Role        CharacterRole `json:"role"`
}

type CreateCharacterDTO struct {
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description" validate:"required"`
	Role        CharacterRole `json:"role" validate:"required"`
}
