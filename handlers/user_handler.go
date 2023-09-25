package handlers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/gorm"

	"github.com/coaltail/GoOrders/database"
	"github.com/coaltail/GoOrders/models"
	"github.com/coaltail/GoOrders/validation"
	"github.com/jinzhu/copier"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var validator = validation.XValidator{}

type WhereFunc func(*gorm.DB) *gorm.DB

func QueryAndReturnError(c *fiber.Ctx, db *gorm.DB, model interface{}, whereFunc WhereFunc) error {
	// Apply the custom search condition using the callback function

	query := whereFunc(db)

	if err := query.First(&model).Error; err != nil {
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
			"message": "Invalid request",
			"error":   err,
		})
	}
	validationErrors := validator.Validate(user)
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validationErrors,
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
	var userProfiles []models.UserProfile
	db := database.DB.Db
	if err := db.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	fmt.Println(users)
	copier.Copy(&userProfiles, &users)
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
		errorMessage := "An error occurred: " + err.Error()
		errorResponse := fiber.Map{"error": errorMessage}

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	var userProfile models.UserProfile
	copier.Copy(&userProfile, &user)
	return c.JSON(userProfile)
}

func UpdateUserProfileByID(c *fiber.Ctx, id string) error {
	var user models.User
	db := database.DB.Db

	err := QueryAndReturnError(c, db, &user, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})

	if err != nil {

		errorMessage := "An error occurred: " + err.Error()
		errorResponse := fiber.Map{"error": errorMessage}

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	// Create a pointer to the user struct
	newUser := &user

	// Parse the request body into newUser
	if err := c.BodyParser(newUser); err != nil {
		return c.JSON(err)
	}
	db.Save(&newUser)
	var userProfile models.UserProfile
	copier.Copy(&userProfile, &newUser)
	return c.JSON(userProfile)
}

func DeleteUserByID(c *fiber.Ctx, id string) error {
	var user models.User
	db := database.DB.Db

	err := QueryAndReturnError(c, db, &user, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})

	if err != nil {
		errorMessage := "An error occurred: " + err.Error()
		errorResponse := fiber.Map{"error": errorMessage}

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete user",
			"errror":  err,
		})
	}

	return c.JSON(fiber.Map{"detail": "success"})
}

func GetUserFollowers(c *fiber.Ctx) error {
	db := database.DB.Db
	fmt.Println("User ID: ", c.Params("id"))
	// Now, retrieve the user's friends using the UserFriends model
	var followers []models.UserFollower
	err := QueryAndReturnError(c, db, &followers, func(db *gorm.DB) *gorm.DB {
		return db.Where("target_id = ?", c.Params("id")).Find(&followers)
	})

	if err != nil {
		errorMessage := "An error occurred: " + err.Error()
		errorResponse := fiber.Map{"error": errorMessage}

		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	return c.JSON(fiber.Map{
		"message":   "Friends retrieved successfully",
		"followers": followers,
	})

}
