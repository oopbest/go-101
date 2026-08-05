package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	_ "github.com/oopbest/go-fiber/model"
	"github.com/oopbest/go-fiber/request"
	"github.com/oopbest/go-fiber/store"
)

var errInvalidID = errors.New("invalid ID")

type BookHandler struct {
	store *store.BookStore
}

func NewBookHandler(store *store.BookStore) *BookHandler {
	return &BookHandler{store: store}
}

// GetAllBooks retrieves all books in the bookstore.
// @Summary      Get all books
// @Description  Get a list of all books stored in memory.
// @Tags         Books
// @Produce      json
// @Success      200  {array}   model.Book
// @Router       /books [get]
func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	return c.JSON(h.store.All())
}

// GetBookByID retrieves a single book by its ID.
// @Summary      Get book by ID
// @Description  Get book details by positive integer ID.
// @Tags         Books
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  model.Book
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /books/{id} [get]
func (h *BookHandler) GetBookByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_ID", "ID must be a positive integer")
	}

	book, found := h.store.FindByID(id)
	if !found {
		return writeError(c, fiber.StatusNotFound, "BOOK_NOT_FOUND", "Book not found")
	}

	return c.JSON(book)
}

// CreateBook adds a new book to the store.
// @Summary      Create a new book
// @Description  Create a book (Requires JWT token with Admin role).
// @Tags         Books
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body request.BookRequest true "Book Data"
// @Success      201  {object}  model.Book
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /books [post]
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	bookRequest, err := parseBookRequest(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_REQUEST", err.Error())
	}

	book := h.store.Create(bookRequest.Title, bookRequest.Author)
	return c.Status(fiber.StatusCreated).JSON(book)
}

// UpdateBook modifies an existing book by its ID.
// @Summary      Update a book
// @Description  Update book details by ID (Requires JWT token with Admin role).
// @Tags         Books
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path      int                 true  "Book ID"
// @Param        request body      request.BookRequest  true  "Updated Book Data"
// @Success      200  {object}  model.Book
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /books/{id} [put]
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_ID", "ID must be a positive integer")
	}

	bookRequest, err := parseBookRequest(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_REQUEST", err.Error())
	}

	book, found := h.store.Update(id, bookRequest.Title, bookRequest.Author)
	if !found {
		return writeError(c, fiber.StatusNotFound, "BOOK_NOT_FOUND", "Book not found")
	}

	return c.JSON(book)
}

// DeleteBook removes a book by its ID.
// @Summary      Delete a book
// @Description  Delete a book by ID (Requires JWT token with Admin role).
// @Tags         Books
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_ID", "ID must be a positive integer")
	}

	if !h.store.Delete(id) {
		return writeError(c, fiber.StatusNotFound, "BOOK_NOT_FOUND", "Book not found")
	}

	return c.JSON(fiber.Map{"message": "Book deleted"})
}

func parseID(c *fiber.Ctx) (int, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return 0, errInvalidID
	}

	return id, nil
}

func parseBookRequest(c *fiber.Ctx) (request.BookRequest, error) {
	var bookRequest request.BookRequest
	if err := c.BodyParser(&bookRequest); err != nil {
		return request.BookRequest{}, errors.New("request body must contain valid JSON")
	}

	bookRequest.Title = strings.TrimSpace(bookRequest.Title)
	bookRequest.Author = strings.TrimSpace(bookRequest.Author)
	if bookRequest.Title == "" {
		return request.BookRequest{}, errors.New("title is required")
	}
	if bookRequest.Author == "" {
		return request.BookRequest{}, errors.New("author is required")
	}

	return bookRequest, nil
}

func writeError(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}
