package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestDatabaseURLFromEnvironment(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("POSTGRES_USER", "postgres")
	t.Setenv("POSTGRES_PASSWORD", "test%#/@: password")
	t.Setenv("POSTGRES_HOST", "postgres")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_DBNAME", "postgres")
	cfg, err := pgx.ParseConfig(databaseURL())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Password != "test%#/@: password" || cfg.Host != "postgres" {
		t.Fatal("environment credentials were changed")
	}
}

// Real route registration with nil dependencies proves rejected requests never
// reach database or storage code. Authorized empty requests reach DTO validation.
func TestAllWriteRoutesRequireAdminToken(t *testing.T) {
	const token = "test-only-admin-token-0123456789abcdef"
	app, err := newApp(nil, nil, token)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/api/characters", "/api/skills", "/api/media-metadata"} {
		for _, tc := range []struct {
			name, header string
			status       int
		}{
			{"anonymous", "", 401},
			{"wrong token", "Bearer wrong", 401},
			{"authorized", "Bearer " + token, 400},
		} {
			t.Run(path+"/"+tc.name, func(t *testing.T) {
				req := httptest.NewRequest("POST", path, strings.NewReader("{}"))
				req.Header.Set("Content-Type", "application/json")
				if tc.header != "" {
					req.Header.Set("Authorization", tc.header)
				}
				res, err := app.Test(req)
				if err != nil {
					t.Fatal(err)
				}
				defer res.Body.Close()
				body, _ := io.ReadAll(res.Body)
				if res.StatusCode != tc.status {
					t.Fatalf("status = %d, want %d; %s", res.StatusCode, tc.status, body)
				}
			})
		}
	}
}

func TestPublicRoutesDoNotRequireAdminToken(t *testing.T) {
	app, err := newApp(nil, nil, "test-only-admin-token-0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/api/characters/bad", "/api/skills/bad", "/api/media-metadata/-1/file"} {
		res, err := app.Test(httptest.NewRequest("GET", path, nil))
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if res.StatusCode != 400 {
			t.Fatalf("public route %s: expected validation error 400, got %d", path, res.StatusCode)
		}
	}
}

func TestAdminEditingRoutesAreProtected(t *testing.T) {
	const token = "test-only-admin-token-0123456789abcdef"
	app, err := newApp(nil, nil, token)
	if err != nil {
		t.Fatal(err)
	}
	for _, route := range []struct {
		method, path     string
		authorizedStatus int
	}{
		{"PUT", "/api/characters/bad", 400},
		{"PUT", "/api/skills/bad", 400},
		{"PUT", "/api/media-metadata/bad", 400},
		{"DELETE", "/api/media-metadata/bad", 400},
		{"POST", "/api/admin/verify", 204},
	} {
		for _, authorized := range []bool{false, true} {
			req := httptest.NewRequest(route.method, route.path, nil)
			status := 401
			if authorized {
				req.Header.Set("Authorization", "Bearer "+token)
				status = route.authorizedStatus
			}
			res, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			_ = res.Body.Close()
			if res.StatusCode != status {
				t.Fatalf("%s %s authorized=%v: got %d, want %d", route.method, route.path, authorized, res.StatusCode, status)
			}
		}
	}
}

func TestExplicitDatabaseURLTakesPriority(t *testing.T) {
	want := "postgres://test:test@localhost:5432/test?sslmode=disable"
	t.Setenv("DATABASE_URL", want)
	if databaseURL() != want {
		t.Fatal("DATABASE_URL should take priority")
	}
}
