package models

import (
	"time"
)

type Book struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null"`
	ISBN        string `gorm:"size:20;uniqueIndex;not null"`
	Price       int64  `gorm:"not null;check:price >= 0"`
	AuthorID    uint   `gorm:"not null;index"`
	Author      Author
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
