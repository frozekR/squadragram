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
	UpdateSkill(context.Context, uuid.UUID, model.CreateSkillDTO) (model.Skill, error)
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
	router.Put("/skills/:uuid", h.UpdateSkill)
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
	return h.saveSkill(c, nil)
}

func (h *Handler) UpdateSkill(c fiber.Ctx) error {
	u, err := httpresponse.UUIDParam(c)
	if err != nil {
		return httpresponse.Error(c, 400, "invalid skill UUID")
	}
	return h.saveSkill(c, &u)
}

func (h *Handler) saveSkill(c fiber.Ctx, u *uuid.UUID) error {
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
	if dto.SortOrder != nil && (*dto.SortOrder < 0 || *dto.SortOrder > 2147483647) {
		return httpresponse.Error(c, 400, "sort_order must be an integer between 0 and 2147483647")
	}
	switch dto.SkillType {
	case model.SkillTypePassive, model.SkillTypeRush, model.SkillTypeSkill,
		model.SkillTypeSuperAttack, model.SkillTypeMaxSuperAttack, model.SkillTypeTransform:
	default:
		return httpresponse.Error(c, fiber.StatusBadRequest, "invalid skill type")
	}
	// CharacterUUID передаётся репозиторию: он найдёт id персонажа для skills.char_id.
	var skill model.Skill
	var err error
	status := fiber.StatusCreated
	if u == nil {
		skill, err = h.repository.CreateSkill(c.Context(), dto)
	} else {
		skill, err = h.repository.UpdateSkill(c.Context(), *u, dto)
		status = fiber.StatusOK
	}
	if err != nil {
		return httpresponse.DatabaseError(c, err, "skill or character not found")
	}
	return c.Status(status).JSON(skill)
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
