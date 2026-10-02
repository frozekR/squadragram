package skill

import (
	"context"
	"squadraton-backend/internal/handler/httpresponse"
	"squadraton-backend/internal/model"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Repository interface {
	CreateSkill(context.Context, model.CreateSkillDTO) (model.Skill, error)
	GetSkillByUUID(context.Context, uuid.UUID) (model.Skill, error)
	GetSkills(context.Context) ([]model.Skill, error)
	GetSkillsByCharacterUUID(context.Context, uuid.UUID) ([]model.Skill, error)
}

type Handler struct {
	repository Repository
}

func NewHandler(repository Repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Get("/skills", h.GetSkills)
	router.Get("/skills/:uuid", h.GetSkillByUUID)
	router.Post("/skills", h.CreateSkill)
	// Вложенный маршрут использует UUID персонажа; репозиторий связывает
	// его characters.id с skills.char_id и возвращает скиллы персонажа.
	router.Get("/characters/:uuid/skills", h.GetSkillsByCharacterUUID)
}

func (h *Handler) GetSkills(c fiber.Ctx) error {
	skills, err := h.repository.GetSkills(c.Context())
	return h.respondSkills(c, skills, err)
}

func (h *Handler) GetSkillByUUID(c fiber.Ctx) error {
	u, err := httpresponse.UUIDParam(c)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid skill UUID")
	}
	skill, err := h.repository.GetSkillByUUID(c.Context(), u)
	if err != nil {
		return httpresponse.DatabaseError(c, err, "skill not found")
	}
	return c.JSON(skill)
}

func (h *Handler) GetSkillsByCharacterUUID(c fiber.Ctx) error {
	u, err := httpresponse.UUIDParam(c)
	if err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid character UUID")
	}
	// Как и в репозитории, неизвестный UUID или отсутствие скиллов дают пустой список.
	skills, err := h.repository.GetSkillsByCharacterUUID(c.Context(), u)
	return h.respondSkills(c, skills, err)
}

func (h *Handler) CreateSkill(c fiber.Ctx) error {
	var dto model.CreateSkillDTO
	if err := c.Bind().JSON(&dto); err != nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid JSON body")
	}
	if !httpresponse.ValidText(dto.Name, dto.Description) {
		return httpresponse.Error(c, fiber.StatusBadRequest, "name (up to 96 characters) and description are required")
	}
	if dto.CharacterUUID == uuid.Nil {
		return httpresponse.Error(c, fiber.StatusBadRequest, "character_uuid is required")
	}
	switch dto.SkillType {
	case model.SkillTypePassive, model.SkillTypeRush, model.SkillTypeSkill,
		model.SkillTypeSuperAttack, model.SkillTypeMaxSuperAttack, model.SkillTypeTransform:
	default:
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid skill type")
	}
	// CharacterUUID передаётся репозиторию: он найдёт id персонажа для skills.char_id.
	skill, err := h.repository.CreateSkill(c.Context(), dto)
	if err != nil {
		return httpresponse.DatabaseError(c, err, "character not found")
	}
	return c.Status(fiber.StatusCreated).JSON(skill)
}

func (h *Handler) respondSkills(c fiber.Ctx, skills []model.Skill, err error) error {
	if err != nil {
		return httpresponse.DatabaseError(c, err, "skill not found")
	}
	if skills == nil {
		skills = []model.Skill{}
	}
	return c.JSON(skills)
}
