package handler

import (
	"crypto/subtle"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oopbest/go-fiber/auth"
	"github.com/oopbest/go-fiber/request"
)

const accessTokenTTL = time.Hour

type AuthHandler struct {
	adminPassword string
	jwtSecret     string
}

func NewAuthHandler(adminPassword, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		adminPassword: adminPassword,
		jwtSecret:     jwtSecret,
	}
}

// Login handles user authentication and generates a JWT token.
// @Summary      User Login
// @Description  Authenticate with username and password to get a JWT token. (Username 'admin' gets admin role, others get user role)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body request.LoginRequest true "Login Credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      429  {object}  map[string]interface{}
// @Router       /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Request body must contain valid JSON")
	}

	username := strings.TrimSpace(loginRequest.Username)
	password := loginRequest.Password
	if username == "" || password == "" {
		return writeError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Username and password are required")
	}

	passwordMatches := subtle.ConstantTimeCompare([]byte(password), []byte(h.adminPassword)) == 1
	if username != "admin" || !passwordMatches {
		return writeError(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password")
	}

	now := time.Now()
	claims := auth.UserClaims{
		Username: username,
		Role:     auth.RoleAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    auth.Issuer,
			Subject:   username,
			Audience:  jwt.ClaimStrings{auth.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return writeError(c, fiber.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate authentication token")
	}

	return c.JSON(fiber.Map{
		"message":    "Login successful",
		"token":      tokenString,
		"token_type": "Bearer",
		"expires_in": int(accessTokenTTL.Seconds()),
	})
}
