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

type WhereFunc func(*gorm.DB) *gorm.DB

func QueryAndReturnError(c *fiber.Ctx, db *gorm.DB, user *models.User, whereFunc WhereFunc) error {
	// Apply the custom search condition using the callback function
	query := whereFunc(db)

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}

		// Handle other database errors
		return errors.New("500 - Internal server error")
	}
	return nil
}

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
	err := QueryAndReturnError(c, db, &user, func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", loginRequest.Email)
	})
	if err != nil {
		return nil
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

	return c.JSON(models.Loginresponse{
		Token: t,
	})

}

func GetUserProfileByID(c *fiber.Ctx, id string) error {
	var user models.User
	db := database.DB.Db

	err := QueryAndReturnError(c, db, &user, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})

	if err != nil {
		return err
	}
	userProfile := models.UserProfile{
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		Mobile:     user.Mobile,
		Email:      user.Email,
		Intro:      user.Intro,
	}
	return c.JSON(userProfile)
}

func UpdateUserProfileByID(c *fiber.Ctx, id string) error {
	var user models.User
	db := database.DB.Db

	err := QueryAndReturnError(c, db, &user, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})
	if err != nil {
		return err
	}

	// Create a pointer to the user struct
	newUser := &user

	// Parse the request body into newUser
	if err := c.BodyParser(newUser); err != nil {
		return err
	}
	db.Save(&newUser)
	return c.JSON(newUser)
}
