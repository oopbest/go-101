package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oopbest/go-fiber/auth"
	"github.com/oopbest/go-fiber/handler"
	"github.com/oopbest/go-fiber/request"
)

func TestAuthHandlerLogin(t *testing.T) {
	const (
		adminPassword = "test-password-123"
		jwtSecret     = "test-jwt-secret-456"
	)

	authHandler := handler.NewAuthHandler(adminPassword, jwtSecret)
	app := fiber.New()
	app.Post("/login", authHandler.Login)

	t.Run("valid admin credentials return an admin JWT", func(t *testing.T) {
		response := login(t, app, request.LoginRequest{
			Username: "admin",
			Password: adminPassword,
		})
		defer response.Body.Close()

		if response.StatusCode != fiber.StatusOK {
			t.Fatalf("expected status %d, got %d", fiber.StatusOK, response.StatusCode)
		}

		var responseBody struct {
			Token     string `json:"token"`
			ExpiresIn int    `json:"expires_in"`
		}
		if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		claims := &auth.UserClaims{}
		token, err := jwt.ParseWithClaims(
			responseBody.Token,
			claims,
			func(_ *jwt.Token) (any, error) { return []byte(jwtSecret), nil },
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)
		if err != nil || !token.Valid {
			t.Fatalf("expected a valid token, got %v", err)
		}
		if claims.Role != auth.RoleAdmin {
			t.Fatalf("expected role %q, got %q", auth.RoleAdmin, claims.Role)
		}
		if claims.Subject != "admin" {
			t.Fatalf("expected subject admin, got %q", claims.Subject)
		}
		if responseBody.ExpiresIn <= 0 {
			t.Fatalf("expected a positive expires_in, got %d", responseBody.ExpiresIn)
		}
	})

	tests := []struct {
		name     string
		request  request.LoginRequest
		expected int
	}{
		{
			name:     "wrong admin password is rejected",
			request:  request.LoginRequest{Username: "admin", Password: "wrong-password"},
			expected: fiber.StatusUnauthorized,
		},
		{
			name:     "unknown user is rejected",
			request:  request.LoginRequest{Username: "john_doe", Password: adminPassword},
			expected: fiber.StatusUnauthorized,
		},
		{
			name:     "empty credentials are rejected",
			request:  request.LoginRequest{},
			expected: fiber.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := login(t, app, test.request)
			defer response.Body.Close()

			if response.StatusCode != test.expected {
				t.Fatalf("expected status %d, got %d", test.expected, response.StatusCode)
			}
		})
	}
}

func login(t *testing.T, app *fiber.App, loginRequest request.LoginRequest) *http.Response {
	t.Helper()

	payload, err := json.Marshal(loginRequest)
	if err != nil {
		t.Fatalf("failed to encode request: %v", err)
	}

	httpRequest := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(payload))
	httpRequest.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	response, err := app.Test(httpRequest)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	return response
}
