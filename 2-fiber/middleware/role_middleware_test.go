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

func TestRequireRoleMiddleware(t *testing.T) {
	const jwtSecret = "my-jwt-secret-key"

	app := fiber.New()
	app.Get(
		"/admin-only",
		middleware.CheckAuth(jwtSecret),
		middleware.RequireRole(auth.RoleAdmin),
		func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) },
	)
	app.Get(
		"/role-without-auth",
		middleware.RequireRole(auth.RoleAdmin),
		func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) },
	)

	tests := []struct {
		name          string
		path          string
		authorization string
		expected      int
	}{
		{
			name:          "admin role is allowed",
			path:          "/admin-only",
			authorization: "Bearer " + createTestJWT(t, jwtSecret, jwt.SigningMethodHS256, auth.RoleAdmin, numericDateAfter(time.Hour)),
			expected:      fiber.StatusOK,
		},
		{
			name:          "user role is forbidden",
			path:          "/admin-only",
			authorization: "Bearer " + createTestJWT(t, jwtSecret, jwt.SigningMethodHS256, auth.RoleUser, numericDateAfter(time.Hour)),
			expected:      fiber.StatusForbidden,
		},
		{
			name:     "role check without authentication is unauthorized",
			path:     "/role-without-auth",
			expected: fiber.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
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
