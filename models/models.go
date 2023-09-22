package models

import (
	"time"
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model

		FirstName    string `gorm:"not null" validate:"required,max=20"`
		MiddleName   string
		LastName     string `gorm:"not null" validate:"required,max=20"`
		Mobile       string `gorm:"unique;not null" validate:"required,min=5,max=20"`
		Email        string `gorm:"unique;not null" validate:"required,min=5,max=45"`
		PasswordHash string `gorm:"not null" validate:"required,min=5,max=85"`
		RegisteredAt time.Time
		LastLogin    time.Time
		Intro        string
	}
)



func AutoMigrate(db *gorm.DB) {
	// AutoMigrate will create the necessary tables in the database
	db.AutoMigrate(&User{})
}
