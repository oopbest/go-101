# การสร้างเอกสาร API ด้วย Swagger (Swagger Documentation in Go & Fiber)

**Swagger (OpenAPI 3.0)** เป็นมาตรฐานการอธิบายสถาปัตยกรรม RESTful API ที่ช่วยให้ผู้พัฒนาสามารถอ่านเอกสารและทดสอบยิง API ผ่านหน้าเว็บเบราว์เซอร์ (Interactive UI) ได้โดยไม่ต้องพึ่งพาคำอธิบายจากภายนอก

---

## 1. แพ็กเกจที่ต้องใช้งาน (Core Packages)

1. **`github.com/gofiber/swagger`**: Fiber Middleware สำหรับให้บริการหน้าเว็บ Swagger UI ผ่าน HTTP Endpoint
2. **`github.com/swaggo/swag`**: เครื่องมือแปลง Annotations ที่เขียนกำกับฟังก์ชัน Go ให้กลายเป็นไฟล์ OpenAPI Specs (`docs/docs.go`, `swagger.json`, `swagger.yaml`)

### การติดตั้ง
```bash
go get github.com/gofiber/swagger
go get github.com/swaggo/swag
```

---

## 2. การระบุ General API Annotations ใน `main.go`

กำหนดข้อมูลพื้นฐานของ API และระบบความปลอดภัย (Security Scheme) ไว้เหนือฟังก์ชัน `main()` ใน [main.go](file:///c:/Users/oopbest-pc-rog/Sites/1_learning/1_golang/mikelopster/2-fiber/main.go):

```go
package main

import (
	"github.com/gofiber/swagger"
	_ "github.com/oopbest/go-fiber/docs" // นำเข้าแพ็กเกจ docs ที่ถูก swag init เจนออกมา
)

// @title           Bookstore Fiber API
// @version         1.0
// @description     RESTful API for Bookstore Management built with Fiber, JWT Auth, and RBAC.
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer ' followed by your JWT token. Example: "Bearer eyJhbGci..."

func main() {
	app := fiber.New()

	// ลงทะเบียน Endpoint สำหรับเปิดหน้าเว็บ Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// ... (routes อื่นๆ) ...
}
```

---

## 3. การเขียน Swag Annotations กำกับใน Handlers

ใส่ Annotations เหนือแต่ละ Handler Function ในแพ็กเกจ `handler`:

### 3.1 ตัวอย่าง: Login Endpoint (`POST /login`)
```go
// Login handles user authentication and generates a JWT token.
// @Summary      User Login
// @Description  Authenticate with username and password to get a JWT token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body request.LoginRequest true "Login Credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error { ... }
```

### 3.2 ตัวอย่าง: Protected Endpoint (`POST /books`)
ใช้ `@Security BearerAuth` เพื่อระบุว่า Endpoint นี้ต้องใช้ JWT Token ในการยืนยันตัวตน:

```go
// CreateBook adds a new book to the store.
// @Summary      Create a new book
// @Description  Create a book (Requires JWT token with Admin role).
// @Tags         Books
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body request.BookRequest true "Book Data"
// @Success      201  {object}  model.Book
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /books [post]
func (h *BookHandler) CreateBook(c *fiber.Ctx) error { ... }
```

### 3.3 ตัวอย่าง: Upload Endpoint (`POST /upload`)
```go
// UploadFile saves an uploaded file to disk.
// @Summary      Upload a document file
// @Tags         Upload
// @Security     BearerAuth
// @Accept       mpfd
// @Produce      json
// @Param        file formData file true "File to upload"
// @Success      201  {object}  map[string]string
// @Router       /upload [post]
func UploadFile(c *fiber.Ctx) error { ... }
```

---

## 4. คำสั่งสร้างไฟล์เอกสาร (Generating OpenAPI Specs)

เมื่อมีการแก้ไขหรือเพิ่ม Annotations ใหม่ ให้สั่งรันคำสั่งเพื่อสร้าง/อัปเดตไฟล์ในโฟลเดอร์ `docs/`:

```bash
go run github.com/swaggo/swag/cmd/swag init
```

ผลลัพธ์ที่จะได้ในไดเรกทอรีโปรเจกต์:
```text
2-fiber/
├── docs/
│   ├── docs.go          # Go source code สำหรับ import ใน main.go
│   ├── swagger.json     # JSON Spec OpenAPI 3.0
│   └── swagger.yaml     # YAML Spec OpenAPI 3.0
```

---

## 5. การเปิดดูและทดสอบยิง API (Interactive UI)

1. รันเซิร์ฟเวอร์ (`go run main.go` หรือ `air`)
2. เปิดเบราว์เซอร์เข้าที่ URL:
   ```text
   http://localhost:8080/swagger/index.html
   ```
3. กดปุ่ม **Authorize** ใส่ `Bearer <jwt_token>` เพื่อทดสอบ Endpoints ที่มีการคุ้มครองสิทธิ์ได้ทันทีผ่านหน้าเว็บ
