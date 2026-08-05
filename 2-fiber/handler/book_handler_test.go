package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oopbest/go-fiber/store"
)

func TestCreateBookRejectsBlankTitle(t *testing.T) {
	bookStore := store.NewBookStore(nil)
	bookHandler := NewBookHandler(bookStore)
	app := fiber.New()
	app.Post("/books", bookHandler.CreateBook)

	httpRequest := httptest.NewRequest(
		http.MethodPost,
		"/books",
		strings.NewReader(`{"title":" ","author":"Author"}`),
	)
	httpRequest.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	response, err := app.Test(httpRequest)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if response.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", fiber.StatusBadRequest, response.StatusCode)
	}
}
