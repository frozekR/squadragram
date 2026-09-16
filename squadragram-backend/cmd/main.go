package main

import (
	"context"
	"log"
	"os"

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

	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":" + envOrDefault("PORT", "3001")))
}

func databaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	return "postgres://" + envOrDefault("POSTGRES_USER", "postgres") + ":" +
		envOrDefault("POSTGRES_PASSWORD", "postgres") + "@" +
		envOrDefault("POSTGRES_HOST", "localhost") + ":" +
		envOrDefault("POSTGRES_PORT", "5432") + "/" +
		envOrDefault("POSTGRES_DBNAME", "postgres") + "?sslmode=disable"
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
