package database

import (
	"context"
	"log"
	"squadraton-backend/internal/config"

	"go.uber.org/fx"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Module = fx.Module("postgres",
	fx.Provide(NewPostgres),
)

func NewPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.Postgres.DSN)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL:", err)
	}

	poolCfg.MinConns = cfg.Postgres.PoolMinConnections
	poolCfg.MaxConns = cfg.Postgres.PoolMaxConnections
	poolCfg.MaxConnIdleTime = cfg.Postgres.PoolMaxIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Fatal("Ошибка создания пула подключений:", err)
	}

	return pool, nil
}
