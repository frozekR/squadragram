package config

import (
	"fmt"
	"time"
)

const (
	PostgresUserEnv     = "POSTGRES_USER"
	PostgresPasswordEnv = "POSTGRES_PASSWORD"
	PostgresHostEnv     = "POSTGRES_HOST"
	PostgresPortEnv     = "POSTGRES_PORT"
	PostgresDBNameEnv   = "POSTGRES_DBNAME"

	defaultPostgresUser     = ""
	defaultPostgresPassword = ""
	defaultPostgresHost     = ""
	defaultPostgresPort     = ""
	defaultPostgresDBName   = ""
)

type Postgres struct {
	DSN string

	PoolMinConnections int32 `yaml:"pool_min_connections"`
	PoolMaxConnections int32 `yaml:"pool_max_connections"`

	PoolMaxIdleTime time.Duration `yaml:"pool_max_idle_time"`
}

func (c *Postgres) loadFromEnv() {
	pgUser := getEnv(PostgresUserEnv, defaultPostgresUser)
	pgPassword := getEnv(PostgresPasswordEnv, defaultPostgresPassword)
	pgHost := getEnv(PostgresHostEnv, defaultPostgresHost)
	pgPort := getEnv(PostgresPortEnv, defaultPostgresPort)
	pgDBName := getEnv(PostgresDBNameEnv, defaultPostgresDBName)

	c.DSN = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgUser,
		pgPassword,
		pgHost,
		pgPort,
		pgDBName,
	)
}
