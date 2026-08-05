package request

// BookRequest contains only the fields that clients are allowed to provide.
type BookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
