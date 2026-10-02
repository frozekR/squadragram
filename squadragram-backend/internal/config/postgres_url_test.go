package config

import (
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestPostgresURLPreservesCredentials(t *testing.T) {
	for _, tc := range []struct{ name, user, password, host, database string }{
		{"ordinary", "postgres", "postgres", "postgres", "postgres"},
		{"special characters", "test@user", "test@:/?#%[]\\ password", "localhost", "test db"},
		{"unicode and IPv6", "пользователь", "тестовый пароль#%", "::1", "test"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := pgx.ParseConfig(PostgresURL(tc.user, tc.password, tc.host, "5432", tc.database))
			if err != nil {
				t.Fatalf("parse connection URL: %v", err)
			}
			if cfg.User != tc.user || cfg.Password != tc.password || cfg.Host != tc.host || cfg.Database != tc.database || cfg.Port != 5432 {
				t.Fatal("connection URL changed the credentials or destination")
			}
		})
	}
}
