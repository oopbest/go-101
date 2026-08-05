package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oopbest/go-fiber/middleware"
)

func TestLoginRateLimit(t *testing.T) {
	app := fiber.New()
	app.Post("/login", middleware.LoginRateLimit(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusUnauthorized)
	})

	for attempt := 1; attempt <= 5; attempt++ {
		request := httptest.NewRequest(http.MethodPost, "/login", nil)
		response, err := app.Test(request)
		if err != nil {
			t.Fatalf("attempt %d failed: %v", attempt, err)
		}
		if response.StatusCode != fiber.StatusUnauthorized {
			t.Fatalf("attempt %d: expected status %d, got %d", attempt, fiber.StatusUnauthorized, response.StatusCode)
		}
	}

	request := httptest.NewRequest(http.MethodPost, "/login", nil)
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("rate-limited request failed: %v", err)
	}
	if response.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d", fiber.StatusTooManyRequests, response.StatusCode)
	}
}
