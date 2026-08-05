package store

import (
	"testing"

	"github.com/oopbest/go-fiber/model"
)

func TestBookStoreCreateKeepsIDUniqueAfterDelete(t *testing.T) {
	bookStore := NewBookStore([]model.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1"},
		{ID: 2, Title: "Book 2", Author: "Author 2"},
	})

	if deleted := bookStore.Delete(1); !deleted {
		t.Fatal("expected book 1 to be deleted")
	}

	created := bookStore.Create("Book 3", "Author 3")
	if created.ID != 3 {
		t.Fatalf("expected the next ID to be 3, got %d", created.ID)
	}
}
