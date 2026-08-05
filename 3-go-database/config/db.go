package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// ConnectDB loads env and connects to PostgreSQL
func ConnectDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, using default environment variables")
	}

	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres123")
	dbname := getEnv("POSTGRES_DB", "godb")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// InitTables creates suppliers and products tables if they don't exist
func InitTables(db *sql.DB) error {
	supplierTableSQL := `
	CREATE TABLE IF NOT EXISTS suppliers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		phone VARCHAR(20),
		email VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	productTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
		stock INT NOT NULL DEFAULT 0,
		supplier_id INT REFERENCES suppliers(id) ON DELETE SET NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(supplierTableSQL); err != nil {
		return fmt.Errorf("failed to create suppliers table: %w", err)
	}

	if _, err := db.Exec(productTableSQL); err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
