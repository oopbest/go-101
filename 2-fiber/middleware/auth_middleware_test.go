package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oopbest/go-fiber/auth"
	"github.com/oopbest/go-fiber/middleware"
)

func TestCheckAuthMiddleware(t *testing.T) {
	const jwtSecret = "my-jwt-secret-key"

	app := fiber.New()
	app.Get("/protected", middleware.CheckAuth(jwtSecret), func(c *fiber.Ctx) error {
		if _, ok := c.Locals("user").(*auth.UserClaims); !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	})

	tests := []struct {
		name          string
		authorization string
		expected      int
	}{
		{name: "missing token", expected: fiber.StatusUnauthorized},
		{
			name:          "valid HS256 token",
			authorization: "Bearer " + createTestJWT(t, jwtSecret, jwt.SigningMethodHS256, auth.RoleAdmin, numericDateAfter(time.Hour)),
			expected:      fiber.StatusOK,
		},
		{
			name:          "expired token",
			authorization: "Bearer " + createTestJWT(t, jwtSecret, jwt.SigningMethodHS256, auth.RoleAdmin, numericDateAfter(-time.Hour)),
			expected:      fiber.StatusUnauthorized,
		},
		{
			name:          "invalid signature",
			authorization: "Bearer " + createTestJWT(t, "wrong-secret", jwt.SigningMethodHS256, auth.RoleAdmin, numericDateAfter(time.Hour)),
			expected:      fiber.StatusUnauthorized,
		},
		{
			name:          "missing expiration",
			authorization: "Bearer " + createTestJWT(t, jwtSecret, jwt.SigningMethodHS256, auth.RoleAdmin, nil),
			expected:      fiber.StatusUnauthorized,
		},
		{
			name:          "HS512 is rejected",
			authorization: "Bearer " + createTestJWT(t, jwtSecret, jwt.SigningMethodHS512, auth.RoleAdmin, numericDateAfter(time.Hour)),
			expected:      fiber.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if test.authorization != "" {
				request.Header.Set(fiber.HeaderAuthorization, test.authorization)
			}

			response, err := app.Test(request)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != test.expected {
				t.Fatalf("expected status %d, got %d", test.expected, response.StatusCode)
			}
		})
	}
}

func createTestJWT(
	t *testing.T,
	secret string,
	method jwt.SigningMethod,
	role auth.Role,
	expiresAt *jwt.NumericDate,
) string {
	t.Helper()

	claims := auth.UserClaims{
		Username: "testuser",
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    auth.Issuer,
			Subject:   "testuser",
			Audience:  jwt.ClaimStrings{auth.Audience},
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(method, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return tokenString
}

func numericDateAfter(duration time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(duration))
}
