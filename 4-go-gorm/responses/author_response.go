package responses

import (
	"time"

	"github.com/oopbest/go-gorm-tut/models"
)

type BookSummaryResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
	Price int64  `json:"price"`
}

type AuthorResponse struct {
	ID        uint                  `json:"id"`
	Name      string                `json:"name"`
	Email     string                `json:"email"`
	Books     []BookSummaryResponse `json:"books"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

func NewAuthorResponse(author models.Author) AuthorResponse {
	books := make(
		[]BookSummaryResponse,
		0,
		len(author.Books),
	)

	for _, book := range author.Books {
		books = append(books, BookSummaryResponse{
			ID:    book.ID,
			Title: book.Title,
			ISBN:  book.ISBN,
			Price: book.Price,
		})
	}

	return AuthorResponse{
		ID:        author.ID,
		Name:      author.Name,
		Email:     author.Email,
		Books:     books,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	}
}

func NewAuthorResponses(
	authors []models.Author,
) []AuthorResponse {
	response := make(
		[]AuthorResponse,
		0,
		len(authors),
	)

	for _, author := range authors {
		response = append(
			response,
			NewAuthorResponse(author),
		)
	}

	return response
}
