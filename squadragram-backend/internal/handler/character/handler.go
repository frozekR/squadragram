package character

import (
	"errors"
	"strconv"

	"squadraton-backend/internal/model"
	characterrepository "squadraton-backend/internal/repository/character"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	repository characterrepository.Querier
}

type createCharacterRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Role        model.CharacterRole `json:"role"`
}

func NewHandler(repository characterrepository.Querier) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Register(router fiber.Router) {
	router.Post("/characters", h.create)
	router.Get("/characters", h.getAll)
	router.Get("/characters/:id", h.get)
}

func (h *Handler) getAll(c fiber.Ctx) error {
	characters, err := h.repository.GetAll(c)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get characters")
	}

	return c.JSON(characters)
}

func (h *Handler) create(c fiber.Ctx) error {
	var request createCharacterRequest
	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if request.Name == "" || !isValidRole(request.Role) {
		return fiber.NewError(fiber.StatusBadRequest, "name and valid role are required")
	}

	character, err := h.repository.Create(c, model.Character{
		Name:        request.Name,
		Description: request.Description,
		Role:        request.Role,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create character")
	}

	return c.Status(fiber.StatusCreated).JSON(character)
}

func (h *Handler) get(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid character id")
	}

	character, err := h.repository.GetByID(c, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return fiber.NewError(fiber.StatusNotFound, "character not found")
	}
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get character")
	}

	return c.JSON(character)
}

func isValidRole(role model.CharacterRole) bool {
	switch role {
	case model.CharacterRoleDamage, model.CharacterRoleTank, model.CharacterRoleTechnical:
		return true
	default:
		return false
	}
}
