package requests

type CreateAuthorRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateAuthorRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
