# การสร้าง RESTful API ด้วย Fiber (Fiber CRUD API)

การพัฒนา **RESTful API** ด้วย **Fiber Framework** ร่วมกับแนวทางการออกแบบโครงสร้างโปรเจกต์ที่เป็นสัดส่วนตามมาตรฐานภาษา Go (Package-based Project Structure) และการจัดการหน่วยความจำแบบพึ่งพาตัวเองได้ (In-memory Storage) ที่รองรับ **Concurrency (Thread-safety)**

---

## 1. โครงสร้างสถาปัตยกรรมโปรเจกต์ (Modular Project Structure)

การปรับปรุงโครงสร้างโปรเจกต์จากไฟล์เดี่ยวไปเป็นแบบแยก **Packages** ช่วยให้โค้ดทดสอบง่ายขึ้น (Testable), บำรุงรักษาง่าย (Maintainable) และขยายขนาดได้ดี (Scalable):

```text
2-fiber/
├── main.go                     # Entry Point ของแอปพลิเคชัน (ตั้งค่า Server & Routes)
├── model/
│   └── book.go                 # Domain Model (package model)
├── request/
│   └── book_request.go         # Request DTOs (package request)
├── store/
│   ├── book_store.go           # In-memory Store & Mutex (package store)
│   └── book_store_test.go      # Unit Test สำหรับ Store
└── handler/
    ├── book_handler.go         # HTTP Handlers / Controllers (package handler)
    └── book_handler_test.go    # Unit Test สำหรับ Handlers
```

---

## 2. โครงสร้างข้อมูลและ DTOs (Models & Requests)

### 2.1 Domain Model ([model/book.go](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/2-fiber/model/book.go))
นิยาม Struct สำหรับข้อมูลหนังสือภายในระบบ:
```go
package model

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
```

### 2.2 Request DTO ([request/book_request.go](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/2-fiber/request/book_request.go))
แยก Struct สำหรับรับ Request Payload จาก Client เพื่อป้องกันไม่ให้ภายนอกระบุ `ID` เข้ามาเอง:
```go
package request

type BookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
```

---

## 3. In-Memory Data Store แบบ Thread-Safe ([store/book_store.go](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/2-fiber/store/book_store.go))

การจัดการข้อมูลในหน่วยความจำโดยใช้ `sync.RWMutex` เพื่อป้องกัน **Data Race** เมื่อมีหลาย Goroutines เข้าถึงพร้อมกัน:

```go
package store

import (
	"sync"
	"github.com/oopbest/go-fiber/model"
)

type BookStore struct {
	mu     sync.RWMutex
	books  []model.Book
	nextID int
}
```

### 3.1 การดึงข้อมูลแบบอ่านเท่านั้น (Read Lock & Defensive Copy)
```go
func (s *BookStore) All() []model.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// คัดลอก (Copy) Slice ใหม่ ป้องกัน Data Race เมื่อถูกนำไปใช้นอก Mutex
	return append([]model.Book(nil), s.books...)
}
```

### 3.2 การเพิ่มและแก้ไขข้อมูล (Write Lock)
```go
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
```

### 3.3 เทคนิคการลบข้อมูลออกจาก Slice (Slice Deletion Trick)
```go
func (s *BookStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, book := range s.books {
		if book.ID == id {
			// ตัดสมาชิกตำแหน่ง i ออก โดยนำส่วนหน้า (0 ถึง i-1) มาต่อกับส่วนหลัง (i+1 ถึงจบ)
			s.books = append(s.books[:i], s.books[i+1:]...)
			return true
		}
	}
	return false
}
```

---

## 4. HTTP Handlers & Controllers ([handler/book_handler.go](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/2-fiber/handler/book_handler.go))

ควบคุมการทำงานผ่าน HTTP Context ของ Fiber โดยรับส่งข้อมูลร่วมกับ `BookStore`:

```go
package handler

type BookHandler struct {
	store *store.BookStore
}

func NewBookHandler(store *store.BookStore) *BookHandler {
	return &BookHandler{store: store}
}
```

### Endpoints Table

| HTTP Method | Path | Handler Method | Status Code | คำอธิบาย |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/books` | `GetAllBooks` | `200 OK` | ดึงข้อมูลหนังสือทั้งหมด |
| `GET` | `/books/:id` | `GetBookByID` | `200 OK` | ดึงข้อมูลหนังสือตาม ID |
| `POST` | `/books` | `CreateBook` | `201 Created` | สร้างหนังสือเล่มใหม่ |
| `PUT` | `/books/:id` | `UpdateBook` | `200 OK` | แก้ไขข้อมูลหนังสือตาม ID |
| `DELETE` | `/books/:id` | `DeleteBook` | `200 OK` | ลบหนังสือตาม ID |

---

## 5. Main Entry Point ([main.go](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/2-fiber/main.go))

ไฟล์ `main.go` ทำหน้าที่เพียงแค่การประกอบ dependencies (Dependency Injection) และตั้งค่า HTTP Router:

```go
package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/oopbest/go-fiber/handler"
	"github.com/oopbest/go-fiber/model"
	"github.com/oopbest/go-fiber/store"
)

func main() {
	bookStore := store.NewBookStore([]model.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1"},
	})
	bookHandler := handler.NewBookHandler(bookStore)

	app := fiber.New()
	app.Get("/books", bookHandler.GetAllBooks)
	app.Get("/books/:id", bookHandler.GetBookByID)
	app.Post("/books", bookHandler.CreateBook)
	app.Put("/books/:id", bookHandler.UpdateBook)
	app.Delete("/books/:id", bookHandler.DeleteBook)

	log.Fatal(app.Listen(":8080"))
}
```

---

## 6. ข้อดีของการแยก Package และทำ Unit Testing

1. **Separation of Concerns**: แต่ละแพ็กเกจรับผิดชอบหน้าที่ของตนเองชัดเจน (`model`, `request`, `store`, `handler`)
2. **Unit Testing**: สามารถเขียนทดสอบ `book_store_test.go` และ `book_handler_test.go` แยกชิ้นส่วนได้อย่างอิสระและครอบคลุม
3. **Clean Main File**: `main.go` เหลือขนาดสั้น สอาด และอ่านเข้าใจง่าย
