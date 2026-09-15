package model

import "github.com/google/uuid"

type Character struct {
	ID          int
	UUID        uuid.UUID
	Name        string
	Description string
	Role        CharacterRole
}
