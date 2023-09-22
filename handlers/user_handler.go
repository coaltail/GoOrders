package handlers

import (
	"errors"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/coaltail/GoOrders/database"
	"github.com/coaltail/GoOrders/models"
	"github.com/coaltail/GoOrders/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	//Hash the password
	password, err := HashPassword(user.PasswordHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
			"error":   err,
		})
	}
	user.PasswordHash = password
	//Validate, to make sure there aren't any empty fields
	validator := validation.XValidator{}
	validationErrors := validator.Validate(user)
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
		})
	}
	db := database.DB.Db
	result := db.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
			"error":   result.Error,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func ListAllUsers(c *fiber.Ctx) error {
	var users []models.User
	db := database.DB.Db
	if err := db.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

func LoginUser(c *fiber.Ctx) error {
	//Extract the login request
	loginRequest := new(models.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//Try to find the user in the database
	var user models.User
	db := database.DB.Db
	if err := db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found in the database
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User not found",
				"error":   err,
			})
		}
		// Other database error occurred
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database error",
			"error":   err,
		})
	}
	//Check the password hash in the database
	if !CheckPasswordHash(loginRequest.Password, user.PasswordHash) {
		// Password doesn't match
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}
	//Make new token
	claims := jwt.MapClaims{
		"ID":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// Delete passwordhash
	user.PasswordHash = ""
	return c.JSON(models.Loginresponse{
		Token: t,
		User:  user,
	})

}
