package character

import (
	"context"

	"squadraton-backend/internal/model"

	"github.com/google/uuid"
)

type Querier interface {
	Create(ctx context.Context, character model.Character) (model.Character, error)
	GetAll(ctx context.Context) ([]model.Character, error)
	GetByID(ctx context.Context, id int64) (model.Character, error)
	GetByUUID(ctx context.Context, id uuid.UUID) (model.Character, error)
}
