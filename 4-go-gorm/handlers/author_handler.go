package handlers

import (
	"errors"
	"log"
	"net/mail"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/oopbest/go-gorm-tut/models"
	"github.com/oopbest/go-gorm-tut/repositories"
	"github.com/oopbest/go-gorm-tut/requests"
	"github.com/oopbest/go-gorm-tut/responses"
	"gorm.io/gorm"
)

type AuthorHandler struct {
	repository *repositories.AuthorRepository
}

func NewAuthorHandler(
	repository *repositories.AuthorRepository,
) *AuthorHandler {
	return &AuthorHandler{
		repository: repository,
	}
}

func (h *AuthorHandler) Create(c fiber.Ctx) error {
	var request requests.CreateAuthorRequest

	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
		)
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.ToLower(
		strings.TrimSpace(request.Email),
	)

	if err := validateAuthorRequest(
		request.Name,
		request.Email,
	); err != nil {
		return err
	}

	author := models.Author{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := h.repository.Create(c.Context(), &author); err != nil {
		return handleRepositoryError(err)
	}

	return c.
		Status(fiber.StatusCreated).
		JSON(responses.NewAuthorResponse(author))
}

func (h *AuthorHandler) FindAll(c fiber.Ctx) error {
	authors, err := h.repository.FindAll(c.Context())
	if err != nil {
		return handleRepositoryError(err)
	}

	return c.JSON(
		responses.NewAuthorResponses(authors),
	)
}

func (h *AuthorHandler) FindByID(c fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid author id",
		)
	}

	author, err := h.repository.FindByID(c.Context(), id)
	if err != nil {
		return handleRepositoryError(err)
	}

	return c.JSON(
		responses.NewAuthorResponse(*author),
	)
}

func (h *AuthorHandler) Update(c fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid author id",
		)
	}

	var request requests.UpdateAuthorRequest

	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid request body",
		)
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.ToLower(
		strings.TrimSpace(request.Email),
	)

	if err := validateAuthorRequest(
		request.Name,
		request.Email,
	); err != nil {
		return err
	}

	if err := h.repository.Update(
		c.Context(),
		id,
		request.Name,
		request.Email,
	); err != nil {
		return handleRepositoryError(err)
	}

	author, err := h.repository.FindByID(c.Context(), id)
	if err != nil {
		return handleRepositoryError(err)
	}

	return c.JSON(
		responses.NewAuthorResponse(*author),
	)
}

func (h *AuthorHandler) Delete(c fiber.Ctx) error {
	id, err := parseID(c.Params("id"))
	if err != nil {
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid author id",
		)
	}

	if err := h.repository.Delete(c.Context(), id); err != nil {
		return handleRepositoryError(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// helpers
func parseID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}

	return uint(id), nil
}

func handleRepositoryError(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return fiber.NewError(
			fiber.StatusNotFound,
			"resource not found",
		)

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return fiber.NewError(
			fiber.StatusConflict,
			"resource already exists",
		)

	case errors.Is(err, gorm.ErrForeignKeyViolated):
		return fiber.NewError(
			fiber.StatusConflict,
			"resource is referenced by another record",
		)

	default:
		log.Printf("repository error: %v", err)

		return fiber.NewError(
			fiber.StatusInternalServerError,
			"internal server error",
		)
	}
}

func isValidEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	if err != nil {
		return false
	}

	return address.Address == value
}

func validateAuthorRequest(name, email string) error {
	switch {
	case name == "" || email == "":
		return fiber.NewError(
			fiber.StatusBadRequest,
			"name and email are required",
		)

	case len([]rune(name)) > 150:
		return fiber.NewError(
			fiber.StatusBadRequest,
			"name must not exceed 150 characters",
		)

	case len([]rune(email)) > 255:
		return fiber.NewError(
			fiber.StatusBadRequest,
			"email must not exceed 255 characters",
		)

	case !isValidEmail(email):
		return fiber.NewError(
			fiber.StatusBadRequest,
			"invalid email format",
		)
	}

	return nil
}
