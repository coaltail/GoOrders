package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Loginresponse struct {
	Token  string     `json:"token"`
	Claims jwt.Claims `json:"claims"`
}

type UserProfile struct {
	ID         uint
	FirstName  string `gorm:"not null" validate:"required,max=20"`
	MiddleName string
	LastName   string `gorm:"not null" validate:"required,max=20"`
	Mobile     string `gorm:"unique;not null" validate:"required,min=5,max=20"`
	Email      string `gorm:"unique;not null" validate:"required,min=5,max=45"`
	Intro      string
}

type QueryFunc func(*gorm.DB) *gorm.DB

// Returns an error if it happens during querying.
func QueryAndReturnError(c *fiber.Ctx, db *gorm.DB, model interface{}, queryFunc QueryFunc) error {
	query := queryFunc(db)

	if err := query.Find(&model).Error; err != nil {
		return err
	}
	return nil
}
