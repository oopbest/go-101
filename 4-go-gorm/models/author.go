package models

import "time"

type Author struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:150;not null"`
	Email     string `gorm:"size:255;uniqueIndex;not null"`
	Books     []Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
