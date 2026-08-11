package repositories

import (
	"context"

	"github.com/oopbest/go-gorm-tut/models"
	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{
		db: db,
	}
}

// Create adds a new author to the database.
func (r *AuthorRepository) Create(
	ctx context.Context,
	author *models.Author,
) error {
	return r.db.
		WithContext(ctx).
		Create(author).
		Error
}

// FindAll retrieves all authors, including their associated books.
func (r *AuthorRepository) FindAll(
	ctx context.Context,
) ([]models.Author, error) {
	var authors []models.Author

	err := r.db.
		WithContext(ctx).
		Preload("Books").
		Order("id ASC").
		Find(&authors).
		Error
	if err != nil {
		return nil, err
	}

	return authors, nil
}

// FindByID retrieves an author by ID, including their associated books.
func (r *AuthorRepository) FindByID(
	ctx context.Context,
	id uint,
) (*models.Author, error) {
	var author models.Author

	err := r.db.
		WithContext(ctx).
		Preload("Books").
		First(&author, id).
		Error
	if err != nil {
		return nil, err
	}

	return &author, nil
}

// Update updates an author's name and email by ID.
func (r *AuthorRepository) Update(
	ctx context.Context,
	id uint,
	name string,
	email string,
) error {
	result := r.db.
		WithContext(ctx).
		Model(&models.Author{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"name":  name,
			"email": email,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete removes an author by ID.
func (r *AuthorRepository) Delete(
	ctx context.Context,
	id uint,
) error {
	result := r.db.
		WithContext(ctx).
		Delete(&models.Author{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
