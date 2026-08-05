# การใช้งาน Go Fiber Framework (Fiber in Go)

**Fiber** เป็น Web Framework สำหรับภาษา Go ที่ถูกออกแบบมาโดยได้รับแรงบันดาลใจจาก **Express.js** (Node.js) ทำงานอยู่บน **Fasthttp** (Engine HTTP ที่เร็วที่สุดใน Go) ทำให้ประมวลผลได้รวดเร็วและใช้หน่วยความจำต่ำมาก เหมาะสำหรับการพัฒนา RESTful API และ Microservices

---

## 1. การติดตั้งและการใช้งานเบื้องต้น (Installation & Basic Setup)

### 1.1 การติดตั้ง Fiber v2
สามารถติดตั้งแพ็กเกจผ่าน `go get`:
```bash
go get -u github.com/gofiber/fiber/v2
```

### 1.2 ตัวอย่าง Server เบื้องต้น
```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// สร้าง Fiber App instance
	app := fiber.New()

	// กำหนด Route GET "/"
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// เริ่มต้น Server ที่ Port 3000
	log.Fatal(app.Listen(":3000"))
}
```

---

## 2. การจัดการ Routing และ Parameters (Routing & Parameters)

### 2.1 HTTP Methods และ Path Parameters
ใช้ `c.Params()` เพื่อดึงค่าจาก URL Path parameter:

```go
// Path Parameter: GET /users/123
app.Get("/users/:id", func(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.SendString("User ID: " + id)
})

// Optional Parameter: GET /hello หรือ GET /hello/john
app.Get("/hello/:name?", func(c *fiber.Ctx) error {
	name := c.Params("name", "Guest") // ถ้าไม่มีค่า ให้ใส่ default เป็น Guest
	return c.SendString("Hello " + name)
})
```

### 2.2 Query Parameters
ใช้ `c.Query()` เพื่ออ่านค่า Query String เช่น `/search?q=golang`:

```go
app.Get("/search", func(c *fiber.Ctx) error {
	searchTerm := c.Query("q", "default-term")
	return c.JSON(fiber.Map{
		"query": searchTerm,
	})
})
```

---

## 3. การรับส่งข้อมูลด้วย JSON (Handling JSON Request & Response)

### 3.1 การส่ง JSON Response
ใช้ `c.JSON()` หรือ `fiber.Map{}` สำหรับส่งคำตอบกลับในรูปแบบ JSON:

```go
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

app.Get("/api/user", func(c *fiber.Ctx) error {
	user := User{ID: 1, Name: "Mikelopster"}
	return c.Status(fiber.StatusOK).JSON(user)
})
```

### 3.2 การรับข้อมูลด้วย BodyParser
ใช้ `c.BodyParser()` เพื่อแปลงข้อมูล JSON ใน Request Body เข้าสู่ Struct:

```go
app.Post("/api/users", func(c *fiber.Ctx) error {
	user := new(User)

	// Bind JSON Body เข้าสู่ Struct
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
})
```

---

## 4. การจัดการ Middleware และ Grouping (Middleware & Route Grouping)

### 4.1 Built-in Middleware (Logger, CORS, Recover)
Fiber มี Middleware ยอดนิยมเตรียมไว้ให้ในตัว:

```go
import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

app.Use(logger.New())  // บันทึก Request Log
app.Use(recover.New()) // ป้องกัน Server ล่มจาก Panic
app.Use(cors.New())    // เปิดใช้งาน CORS
```

### 4.2 Route Grouping (จัดกลุ่ม API)
จัดหมวดหมู่ Endpoint เช่น `/api/v1`:

```go
api := app.Group("/api/v1")

api.Get("/products", getProducts)
api.Post("/products", createProduct)
```

---

## 5. สรุปเปรียบเทียบ `net/http` vs `Fiber`

| คุณสมบัติ | `net/http` (Standard Library) | `Fiber` Framework |
| :--- | :--- | :--- |
| **ความสะดวกในการ Routing** | ต้องจัดการ Pattern เอง | มี Router สำเร็จรูปสไตล์ Express.js |
| **JSON Response** | ต้องใช้องค์ประกอบ `json.Encoder` | มี `c.JSON()` และ `c.BodyParser()` |
| **Middleware** | เขียนแบบ Wrapper Pattern เอง | มี Middleware ชนิดพร้อมใช้งานหลากหลาย |
| **Performance** | สูงมาก (Standard Go) | สูงมากเป็นพิเศษ (ขับเคลื่อนด้วย Fasthttp) |
