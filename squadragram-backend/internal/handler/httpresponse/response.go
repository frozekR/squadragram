package httpresponse

import (
	"errors"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func Error(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}

// DatabaseError keeps driver details out of HTTP responses while preserving useful statuses.
func DatabaseError(c fiber.Ctx, err error, notFound string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return Error(c, fiber.StatusNotFound, notFound)
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return Error(c, fiber.StatusConflict, "name already exists")
		case "23503":
			return Error(c, fiber.StatusNotFound, "character not found")
		}
	}
	log.Printf("database request failed: %v", err)
	return Error(c, fiber.StatusInternalServerError, "internal server error")
}

func UUIDParam(c fiber.Ctx) (uuid.UUID, error) {
	u, err := uuid.Parse(c.Params("uuid"))
	if err != nil || u == uuid.Nil {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, "invalid UUID")
	}
	return u, nil
}

func ValidText(name, description string) bool {
	return strings.TrimSpace(name) != "" && utf8.RuneCountInString(name) <= 96 &&
		strings.TrimSpace(description) != ""
}
