package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oopbest/go-gorm-tut/handlers"
)

func RegisterAuthorRoutes(
	router fiber.Router,
	handler *handlers.AuthorHandler,
) {
	router.Post("/authors", handler.Create)
	router.Get("/authors", handler.FindAll)
	router.Get("/authors/:id", handler.FindByID)
	router.Put("/authors/:id", handler.Update)
	router.Delete("/authors/:id", handler.Delete)
}
