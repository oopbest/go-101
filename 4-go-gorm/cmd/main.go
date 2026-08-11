package main

import (
	"log"

	"github.com/oopbest/go-gorm-tut/config"
	"github.com/oopbest/go-gorm-tut/models"
)

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	if err := db.AutoMigrate(
		&models.Author{},
		&models.Book{},
	); err != nil {
		log.Fatal(err)
	}

	log.Println("database connected and migrated successfully")
}
