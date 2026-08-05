package auth

import "github.com/golang-jwt/jwt/v5"

const (
	Issuer   = "go-fiber-api"
	Audience = "go-fiber-client"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// UserClaims contains the authenticated identity and standard JWT claims.
type UserClaims struct {
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}
