package character

import (
	"context"
	"squadraton-backend/internal/handler/httpresponse"
	"squadraton-backend/internal/model"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Repository interface {
	CreateCharacter(context.Context, model.CreateCharacterDTO) (model.Character, error)
	UpdateCharacter(context.Context, uuid.UUID, model.CreateCharacterDTO) (model.Character, error)
	GetCharacterByUUID(context.Context, uuid.UUID) (model.Character, error)
	GetCharacters(context.Context) ([]model.Character, error)
}

type Handler struct {
	repository Repository
}

func NewHandler(repository Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/characters", h.GetCharacters)
	router.Get("/characters/:uuid", h.GetCharacterByUUID)
	router.Post("/characters", h.CreateCharacter)
	router.Put("/characters/:uuid", h.UpdateCharacter)
}

func (h *Handler) GetCharacters(c fiber.Ctx) error {
	characters, err := h.repository.GetCharacters(c.Context())
	if err != nil {
		return httpresponse.DatabaseError(c, err, "character not found")
	}
	if characters == nil {
		characters = []model.Character{}
	}
	return c.JSON(characters)
}

func (h *Handler) GetCharacterByUUID(c fiber.Ctx) error {
	u, err := httpresponse.UUIDParam(c)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid character UUID")
	}
	character, err := h.repository.GetCharacterByUUID(c.Context(), u)
	if err != nil {
		return httpresponse.DatabaseError(c, err, "character not found")
	}
	return c.JSON(character)
}

func (h *Handler) CreateCharacter(c fiber.Ctx) error {
	return h.saveCharacter(c, nil)
}

func (h *Handler) UpdateCharacter(c fiber.Ctx) error {
	u, err := httpresponse.UUIDParam(c)
	if err != nil {
		return httpresponse.Error(c, 400, "invalid character UUID")
	}
	return h.saveCharacter(c, &u)
}

func (h *Handler) saveCharacter(c fiber.Ctx, u *uuid.UUID) error {
	var dto model.CreateCharacterDTO
	if err := c.Bind().JSON(&dto); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid JSON body")
	}
	if !httpresponse.ValidText(dto.Name, dto.Description) {
		return httpresponse.Error(c, fiber.StatusBadRequest, "name (up to 96 characters) and description are required")
	}
	switch dto.Role {
	case model.CharacterRoleDamage, model.CharacterRoleTank, model.CharacterRoleTechnical,
		model.CharacterRoleMelee, model.CharacterRoleRanged:
	default:
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid character role")
	}
	var character model.Character
	var err error
	status := fiber.StatusCreated
	if u == nil {
		character, err = h.repository.CreateCharacter(c.Context(), dto)
	} else {
		character, err = h.repository.UpdateCharacter(c.Context(), *u, dto)
		status = fiber.StatusOK
	}
	if err != nil {
		return httpresponse.DatabaseError(c, err, "character not found")
	}
	return c.Status(status).JSON(character)
}
