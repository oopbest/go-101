package request

// LoginRequest contains the credentials provided by the client.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
