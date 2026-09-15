package model

import (
	"github.com/google/uuid"
)

type Skill struct {
	ID          int
	UUID        uuid.UUID
	Name        string
	Description string
	SkillType   SkillType
}
