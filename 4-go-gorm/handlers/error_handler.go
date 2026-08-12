package handlers

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "internal server error"

	var fiberErr *fiber.Error

	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	} else {
		log.Printf("unhandled error: %v", err)
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}
