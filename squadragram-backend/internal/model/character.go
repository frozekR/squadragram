package model

import "github.com/google/uuid"

type Character struct {
	ID          int           `json:"id"`
	UUID        uuid.UUID     `json:"uuid"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Role        CharacterRole `json:"role"`
}
