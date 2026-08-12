package main

import (
	"context"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recoverer "github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/oopbest/go-gorm-tut/config"
	"github.com/oopbest/go-gorm-tut/handlers"
	"github.com/oopbest/go-gorm-tut/models"
	"github.com/oopbest/go-gorm-tut/repositories"
	"github.com/oopbest/go-gorm-tut/routes"
)

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := db.AutoMigrate(
		&models.Author{},
		&models.Book{},
	); err != nil {
		log.Fatal(err)
	}
	log.Println("database connected and migrated successfully")

	authorRepository := repositories.NewAuthorRepository(db)
	authorHandler := handlers.NewAuthorHandler(authorRepository)

	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.ErrorHandler,
	})

	var shuttingDown atomic.Bool

	app.Use(recoverer.New())
	app.Use(logger.New())

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/ready", func(c fiber.Ctx) error {
		if shuttingDown.Load() {
			return c.
				Status(fiber.StatusServiceUnavailable).
				JSON(fiber.Map{
					"status": "not_ready",
				})
		}

		ctx, cancel := context.WithTimeout(
			c,
			2*time.Second,
		)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			return c.
				Status(fiber.StatusServiceUnavailable).
				JSON(fiber.Map{
					"status":   "not_ready",
					"database": "unavailable",
				})
		}

		return c.JSON(fiber.Map{
			"status":   "ready",
			"database": "ok",
		})
	})

	api := app.Group("/api/v1")

	routes.RegisterAuthorRoutes(
		api,
		authorHandler,
	)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Printf("server stopped: %v", err)
	}

}
