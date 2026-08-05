package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oopbest/go-fiber/auth"
)

const userContextKey = "user"

// CheckAuth verifies a Bearer JWT and stores its typed claims in the request context.
func CheckAuth(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := bearerToken(c.Get(fiber.HeaderAuthorization))
		if tokenString == "" {
			return unauthorized(c, "Missing authentication token")
		}

		claims := &auth.UserClaims{}
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(_ *jwt.Token) (any, error) {
				return []byte(jwtSecret), nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			jwt.WithExpirationRequired(),
			jwt.WithIssuer(auth.Issuer),
			jwt.WithAudience(auth.Audience),
			jwt.WithIssuedAt(),
		)
		if err != nil || token == nil || !token.Valid {
			return unauthorized(c, "Invalid or expired JWT token")
		}

		parsedClaims, ok := token.Claims.(*auth.UserClaims)
		if !ok {
			return unauthorized(c, "Invalid JWT claims")
		}

		c.Locals(userContextKey, parsedClaims)
		return c.Next()
	}
}

func bearerToken(authorization string) string {
	scheme, token, found := strings.Cut(strings.TrimSpace(authorization), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}

	return strings.TrimSpace(token)
}

func unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
	})
}
