package handler

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

// UploadFile saves an uploaded file to disk.
// @Summary      Upload a document file
// @Description  Upload file via multipart/form-data (Requires JWT token with Admin role).
// @Tags         Upload
// @Security     BearerAuth
// @Accept       mpfd
// @Produce      json
// @Param        file formData file true "File to upload"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /upload [post]
func UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	if err := os.MkdirAll("uploads", 0o755); err != nil {
		return err
	}

	filename := filepath.Base(file.Filename)
	if err := c.SaveFile(file, filepath.Join("uploads", filename)); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"filename": filename})
}
