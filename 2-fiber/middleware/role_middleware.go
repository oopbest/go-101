package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oopbest/go-fiber/auth"
)

// RequireRole allows only authenticated users with the required role.
func RequireRole(requiredRole auth.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(userContextKey).(*auth.UserClaims)
		if !ok {
			return unauthorized(c, "Authentication required before role check")
		}

		if claims.Role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fiber.Map{
					"code":    "FORBIDDEN",
					"message": "Access denied: insufficient permissions",
				},
			})
		}

		return c.Next()
	}
}
