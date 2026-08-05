package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"github.com/oopbest/go-fiber/auth"
	_ "github.com/oopbest/go-fiber/docs"
	"github.com/oopbest/go-fiber/handler"
	"github.com/oopbest/go-fiber/middleware"
	"github.com/oopbest/go-fiber/model"
	"github.com/oopbest/go-fiber/store"
)

// @title           Bookstore Fiber API
// @version         1.0
// @description     RESTful API for Bookstore Management built with Fiber, JWT Auth, and Role-Based Access Control.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer ' followed by your JWT token. Example: "Bearer eyJhbGci..."

func getEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	log.Fatalf("Environment variable %s is required but not set", key)
	return ""
}

func main() {
	// A .env file is optional when variables are injected by the environment.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Error loading .env file: %v", err)
	}
	adminPassword := getEnv("ADMIN_PASSWORD")
	jwtSecret := getEnv("JWT_SECRET", "defaultjwtsecretkey")

	bookStore := store.NewBookStore([]model.Book{
		{ID: 1, Title: "Book 1", Author: "Author 1"},
		{ID: 2, Title: "Book 2", Author: "Author 2"},
	})
	bookHandler := handler.NewBookHandler(bookStore)
	authHandler := handler.NewAuthHandler(adminPassword, jwtSecret)

	app := fiber.New()

	// Swagger UI Route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Public Routes
	app.Post("/login", middleware.LoginRateLimit(), authHandler.Login)
	app.Get("/books", bookHandler.GetAllBooks)
	app.Get("/books/:id", bookHandler.GetBookByID)

	// Protected Admin Routes (Require Authentication with JWT & Admin Role)
	protected := app.Group("/", middleware.CheckAuth(jwtSecret), middleware.RequireRole(auth.RoleAdmin))
	protected.Post("/books", bookHandler.CreateBook)
	protected.Put("/books/:id", bookHandler.UpdateBook)
	protected.Delete("/books/:id", bookHandler.DeleteBook)
	protected.Post("/upload", handler.UploadFile)

	log.Fatal(app.Listen(":8080"))
}
