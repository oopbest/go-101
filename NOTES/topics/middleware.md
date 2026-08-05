# การใช้งาน Middleware ในภาษา Go และ Fiber (Middleware in Go & Fiber)

**Middleware** คือ ฟังก์ชันที่ทำหน้าที่เป็น "ตัวกลาง" แทรกอยู่ระหว่างที่ **HTTP Request เข้ามา** จนกระทั่งส่ง **HTTP Response กลับออกไป** ใน Request-Response Pipeline ใช้สำหรับประมวลผลงานส่วนกลางที่ทุก Endpoint ต้องใช้ร่วมกัน เช่น การบันทึก Log, การทำ Authentication, การดักจับ Panic หรือการกำหนด CORS

---

## 1. ลำดับการทำงานของ Middleware (Middleware Pipeline)

```text
[HTTP Request] 
      │
      ▼
┌──────────────┐
│  Middleware1 │ (เช่น Logger / CORS)
└──────┬───────┘
       │ Next()
       ▼
┌──────────────┐
│  Middleware2 │ (เช่น Auth / Secret Check)
└──────┬───────┘
       │ Next()
       ▼
┌──────────────┐
│  Handler/API │ (ประมวลผล Business Logic จริง)
└──────┬───────┘
       │ Return
       ▼
[HTTP Response]
```

---

## 2. Built-in Middlewares ใน Fiber Framework

Fiber มี Middleware ยอดนิยมเตรียมมาให้พร้อมใช้งานในตัว สามารถเรียกใช้ผ่าน `app.Use()`:

### 2.1 Logger Middleware
เก็บบันทึกข้อมูล HTTP Request แต่ละรายการ (Method, Path, Status Code, Latency):

```go
import "github.com/gofiber/fiber/v2/middleware/logger"

app := fiber.New()
app.Use(logger.New()) // บันทึก Log ทุกครั้งที่มี Request เข้ามา
```

### 2.2 Recover Middleware
ดักจับข้อผิดพลาดรุนแรงระดับ **Panic** ไม่ให้ส่งผลให้ Web Server ล่ม (Crash) และเปลี่ยนให้ส่ง HTTP Status `500 Internal Server Error` แทน:

```go
import "github.com/gofiber/fiber/v2/middleware/recover"

app.Use(recover.New())
```

### 2.3 CORS Middleware
จัดการสิทธิ์การเข้าถึง API ข้ามโดเมน (Cross-Origin Resource Sharing):

```go
import "github.com/gofiber/fiber/v2/middleware/cors"

app.Use(cors.New(cors.Config{
	AllowOrigins: "https://example.com, https://myfrontend.com",
	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
}))
```

---

## 3. การเขียน Custom Middleware ขึ้นใช้งานเอง (Custom Middleware)

ฟังก์ชัน Middleware ใน Fiber จะรับพารามิเตอร์ `c *fiber.Ctx` และต้องคืนค่าเป็น `error` หากต้องการส่งคำขอไปยัง Handler ถัดไป ต้องเรียก `c.Next()`

### 3.1 ตัวอย่าง: Secret API Key Check Middleware
สร้าง Middleware สำหรับตรวจสอบว่า Request ที่ส่งมามี `X-Secret-Key` ตรงตามที่กำหนดหรือไม่:

```go
package main

import (
	"os"
	"github.com/gofiber/fiber/v2"
)

func SecretAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		secretKey := c.Get("X-Secret-Key")
		expectedSecret := os.Getenv("SECRET")

		if secretKey == "" || secretKey != expectedSecret {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid or missing Secret Key",
			})
		}

		// ผ่านการตรวจสอบ ส่งต่อไปยัง Handler ถัดไป
		return c.Next()
	}
}
```

### 3.2 ตัวอย่าง: Request Duration / Timing Middleware
คำนวณเวลาที่ใช้ในการประมวลผล Request แต่ละครั้ง:

```go
func RequestTimer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// ส่งให้ Handler ประมวลผลก่อน
		err := c.Next()

		duration := time.Since(start)
		log.Printf("[%s] %s took %v", c.Method(), c.Path(), duration)

		return err
	}
}
```

### 3.3 ตัวอย่าง: Role-Based Access Control (RBAC) Middleware
ตรวจสอบว่าผู้ใช้งานที่ยิง Request เข้ามามีสิทธิ์ (Role) ตรงตามที่กำหนดหรือไม่ (เช่น ต้องเป็น `"admin"` เท่านั้น):

```go
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userLocals := c.Locals("user")
		if userLocals == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		claims, ok := userLocals.(jwt.MapClaims)
		if !ok || claims["role"] != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: insufficient permissions",
			})
		}

		return c.Next()
	}
}
```

---

## 4. การประยุกต์ใช้งานร่วมกับ Route Grouping

สามารถเลือกให้ Middleware ทำงานเฉพาะกับกลุ่ม API ที่กำหนดได้โดยใช้ `app.Group()`:

```go
app := fiber.New()

// Public Endpoints (ไม่ต้องยืนยันตัวตน)
app.Get("/public", func(c *fiber.Ctx) error {
	return c.SendString("Public Content")
})

// Protected Endpoints Group (ใช้ SecretAuthMiddleware)
api := app.Group("/api/v1", SecretAuthMiddleware())

api.Get("/books", bookHandler.GetAllBooks)
api.Post("/books", bookHandler.CreateBook)
```

---

## 5. สรุปประโยชน์ของ Middleware

1. **Reusability**: เขียน Logic สำหรับตรวจสอบหรือบันทึกข้อมูลเพียงครั้งเดียว แต่ใช้ร่วมกันได้ทุก Endpoint
2. **Clean Code (Separation of Concerns)**: แยกโค้ดตรวจสอบสิทธิ์ (Auth/Validation) ออกจาก Business Logic ของ Controller/Handler
3. **Security & Reliability**: ช่วยเพิ่มความปลอดภัย (CORS, Auth) และความเสถียร (Recover จาก Panic) ให้กับระบบ API
