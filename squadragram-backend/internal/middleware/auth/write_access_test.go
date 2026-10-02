package auth

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

const testToken = "test-only-admin-token-0123456789abcdef"

func TestWriteAccess(t *testing.T) {
	for _, tc := range []struct {
		name, method, header, query, token string
		status                             int
	}{
		{"public GET", "GET", "", "", testToken, 204},
		{"public HEAD", "HEAD", "", "", testToken, 204},
		{"preflight", "OPTIONS", "", "", testToken, 204},
		{"missing POST token", "POST", "", "", testToken, 401},
		{"missing PUT token", "PUT", "", "", testToken, 401},
		{"missing PATCH token", "PATCH", "", "", testToken, 401},
		{"missing DELETE token", "DELETE", "", "", testToken, 401},
		{"wrong token", "POST", "Bearer wrong", "", testToken, 401},
		{"wrong scheme", "POST", "Basic " + testToken, "", testToken, 401},
		{"extra credentials", "POST", "Bearer " + testToken + " extra", "", testToken, 401},
		{"URL token rejected", "POST", "", "?token=" + testToken, testToken, 401},
		{"valid token", "POST", "Bearer " + testToken, "", testToken, 204},
		{"case insensitive scheme", "DELETE", "bearer " + testToken, "", testToken, 204},
		{"unconfigured writes", "POST", "", "", "", 503},
		{"unconfigured with token", "POST", "Bearer " + testToken, "", "", 503},
		{"unconfigured reads", "GET", "", "", "", 204},
	} {
		t.Run(tc.name, func(t *testing.T) {
			guard, err := WriteAccess(tc.token)
			if err != nil {
				t.Fatal(err)
			}
			app := fiber.New()
			called := false
			app.Use("/api", guard)
			app.All("/api/test", func(c fiber.Ctx) error { called = true; return c.SendStatus(204) })
			req := httptest.NewRequest(tc.method, "/api/test"+tc.query, nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			res, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != tc.status {
				t.Fatalf("status = %d, want %d", res.StatusCode, tc.status)
			}
			if called != (tc.status == 204) {
				t.Fatal("unexpected handler access")
			}
			if tc.status == 401 && res.Header.Get("WWW-Authenticate") == "" {
				t.Fatal("missing authentication challenge")
			}
			if tc.status != 204 && res.Header.Get("Cache-Control") != "no-store" {
				t.Fatal("auth response must not be cached")
			}
		})
	}
}

func TestRejectWeakOrInvalidConfiguration(t *testing.T) {
	for _, token := range []string{"short", strings.Repeat("x", 31), strings.Repeat("x", 32) + "\n", strings.Repeat("x", 32) + " "} {
		if _, err := WriteAccess(token); err == nil {
			t.Fatal("expected invalid token configuration to fail")
		}
	}
}
