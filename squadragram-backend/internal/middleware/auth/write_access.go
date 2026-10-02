package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"strings"

	"squadraton-backend/internal/handler/httpresponse"

	"github.com/gofiber/fiber/v3"
)

// WriteAccess protects all mutating API methods, including routes added later.
// An empty token disables writes rather than making them public.
func WriteAccess(token string) (fiber.Handler, error) {
	if token != "" && (len(token) < 32 || strings.ContainsAny(token, " \t\r\n")) {
		return nil, errors.New("ADMIN_API_TOKEN must contain at least 32 bytes and no whitespace")
	}
	expected := sha256.Sum256([]byte(token))
	return func(c fiber.Ctx) error {
		switch c.Method() {
		case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
			return c.Next()
		}
		c.Set("Cache-Control", "no-store")
		if token == "" {
			return httpresponse.Error(c, fiber.StatusServiceUnavailable, "API writes are disabled")
		}
		// Credentials are accepted only in the Authorization header, never a URL.
		parts := strings.Fields(c.Get("Authorization"))
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			provided := sha256.Sum256([]byte(parts[1]))
			if subtle.ConstantTimeCompare(provided[:], expected[:]) == 1 {
				return c.Next()
			}
		}
		c.Set("WWW-Authenticate", `Bearer realm="admin"`)
		return httpresponse.Error(c, fiber.StatusUnauthorized, "unauthorized")
	}, nil
}
