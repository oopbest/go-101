package repositories

import (
	"context"

	"github.com/oopbest/go-gorm-tut/models"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

// Create adds a new book to the database.
func (r *BookRepository) Create(
	ctx context.Context,
	book *models.Book,
) error {
	return r.db.
		WithContext(ctx).
		Create(book).
		Error
}

// FindAll retrieves all books from the database, including their associated authors.
func (r *BookRepository) FindAll(
	ctx context.Context,
) ([]models.Book, error) {
	var books []models.Book

	err := r.db.
		WithContext(ctx).
		Preload("Author").
		Order("id ASC").
		Find(&books).
		Error
	if err != nil {
		return nil, err
	}

	return books, nil
}

// FindByID retrieves a book by its ID from the database, including its associated author.
func (r *BookRepository) FindByID(
	ctx context.Context,
	id uint,
) (*models.Book, error) {
	var book models.Book

	err := r.db.
		WithContext(ctx).
		Preload("Author").
		First(&book, id).
		Error
	if err != nil {
		return nil, err
	}

	return &book, nil
}

// Update updates a book's title, ISBN, price, author ID, and published date by ID.
func (r *BookRepository) Update(
	ctx context.Context,
	book *models.Book,
) error {
	result := r.db.
		WithContext(ctx).
		Model(&models.Book{}).
		Where("id = ?", book.ID).
		Updates(map[string]any{
			"title":        book.Title,
			"isbn":         book.ISBN,
			"price":        book.Price,
			"author_id":    book.AuthorID,
			"published_at": book.PublishedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete removes a book from the database by its ID.
func (r *BookRepository) Delete(
	ctx context.Context,
	id uint,
) error {
	result := r.db.
		WithContext(ctx).
		Delete(&models.Book{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
