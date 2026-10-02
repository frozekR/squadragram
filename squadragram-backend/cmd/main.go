package main

import (
	"context"
	"log"
	"os"
	"squadraton-backend/internal/config"
	characterhandler "squadraton-backend/internal/handler/character"
	mediahandler "squadraton-backend/internal/handler/media_metadata"
	skillhandler "squadraton-backend/internal/handler/skill"
	"squadraton-backend/internal/middleware/auth"
	"squadraton-backend/internal/repository"
	characterrepository "squadraton-backend/internal/repository/character"
	mediarepository "squadraton-backend/internal/repository/media_metadata"
	skillrepository "squadraton-backend/internal/repository/skill"
	storage "squadraton-backend/pkg/minio"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	db, err := pgxpool.New(context.Background(), databaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mediaStorage, err := storage.NewMinio()
	if err != nil {
		log.Fatal(err)
	}
	app, err := newApp(db, mediaStorage, os.Getenv("ADMIN_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(app.Listen(":" + envOrDefault("PORT", "3001")))
}

func newApp(db repository.DB, mediaStorage mediahandler.Storage, adminToken string) (*fiber.App, error) {
	writeAccess, err := auth.WriteAccess(adminToken)
	if err != nil {
		return nil, err
	}
	app := fiber.New(fiber.Config{BodyLimit: int(mediahandler.MaxFileSize) + (1 << 20)})
	app.Use(cors.New())
	api := app.Group("/api")
	// Register before every API route, so no write handler can bypass authorization.
	api.Use(writeAccess)
	api.Post("/admin/verify", func(c fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })
	characterhandler.NewHandler(characterrepository.NewRepository(db)).RegisterRoutes(api)
	skillhandler.NewHandler(skillrepository.NewRepository(db)).RegisterRoutes(api)
	mediahandler.NewHandler(mediarepository.NewRepository(db), mediaStorage).RegisterRoutes(api)

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	return app, nil
}

func databaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	return config.PostgresURL(
		envOrDefault("POSTGRES_USER", "postgres"),
		envOrDefault("POSTGRES_PASSWORD", "postgres"),
		envOrDefault("POSTGRES_HOST", "postgres"),
		envOrDefault("POSTGRES_PORT", "5432"),
		envOrDefault("POSTGRES_DBNAME", "postgres"),
	)
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
