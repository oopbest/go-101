package store

import (
	"sync"

	"github.com/oopbest/go-fiber/model"
)

// BookStore owns the in-memory data and protects it from concurrent access.
type BookStore struct {
	mu     sync.RWMutex
	books  []model.Book
	nextID int
}

func NewBookStore(seed []model.Book) *BookStore {
	bookStore := &BookStore{
		books:  append([]model.Book(nil), seed...),
		nextID: 1,
	}

	for _, book := range seed {
		if book.ID >= bookStore.nextID {
			bookStore.nextID = book.ID + 1
		}
	}

	return bookStore
}

func (s *BookStore) All() []model.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]model.Book(nil), s.books...)
}

func (s *BookStore) FindByID(id int) (model.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, book := range s.books {
		if book.ID == id {
			return book, true
		}
	}

	return model.Book{}, false
}

func (s *BookStore) Create(title, author string) model.Book {
	s.mu.Lock()
	defer s.mu.Unlock()

	book := model.Book{
		ID:     s.nextID,
		Title:  title,
		Author: author,
	}
	s.nextID++
	s.books = append(s.books, book)

	return book
}

func (s *BookStore) Update(id int, title, author string) (model.Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.books {
		if s.books[i].ID == id {
			updatedBook := model.Book{
				ID:     id,
				Title:  title,
				Author: author,
			}
			s.books[i] = updatedBook
			return updatedBook, true
		}
	}

	return model.Book{}, false
}

func (s *BookStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, book := range s.books {
		if book.ID == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			return true
		}
	}

	return false
}
