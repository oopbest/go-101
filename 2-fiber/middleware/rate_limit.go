package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

const maxLoginAttempts = 5

// LoginRateLimit limits failed login attempts per IP address.
func LoginRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:                    maxLoginAttempts,
		Expiration:             time.Minute,
		SkipSuccessfulRequests: true,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many login attempts. Please try again later",
				},
			})
		},
	})
}
