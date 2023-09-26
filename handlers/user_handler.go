package handlers

import (
	"errors"
	"os"
	"strconv"
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

// Helper functions and interfaces
var validator = validation.XValidator{}

func CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"error":   err,
		})
	}
	validation_errors := validator.Validate(user)
	if len(validation_errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": validation_errors,
		})
	}
	//Hash the password
	password, err := HashPassword(user.PasswordHash)
	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Failed to hash password", err)
	}
	user.PasswordHash = password
	db := database.DB.Db
	result := db.Create(&user)
	if result.Error != nil {
		return handleError(c, fiber.StatusInternalServerError, "Failed to create user", result.Error)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func ListAllUsers(c *fiber.Ctx) error {
	var users []models.User
	var userProfiles []models.UserProfile
	db := database.DB.Db
	if err := db.Find(&users).Error; err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not find users", err)
	}
	copier.Copy(&userProfiles, &users)
	return c.JSON(userProfiles)
}

func LoginUser(c *fiber.Ctx) error {
	// Extract the login request
	loginRequest := new(models.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return handleError(c, fiber.StatusBadRequest, "Could not parse request", err)
	}

	// Try to find the user in the database
	var user models.User
	db := database.DB.Db
	err := models.QueryAndReturnError(c, db, &user, func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", loginRequest.Email)
	})
	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Failed to log in user", err)
	}

	// Check the password hash in the database
	if !CheckPasswordHash(loginRequest.Password, user.PasswordHash) {
		// Password doesn't match
		return handleError(c, fiber.StatusUnauthorized, "Invalid credentials", err)
	}

	// Check if the user already has a token
	var existingToken models.Token
	tokenDB := db.Model(&user).Association("Token")
	tokenDB.Find(&existingToken)

	// Generate or update the token
	tokenExpiry := time.Now().Add(time.Hour * 72).Unix()
	claims := jwt.MapClaims{
		"ID":    user.ID,
		"email": user.Email,
		"exp":   tokenExpiry,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Error signing token", err)
	}

	if existingToken.ID == 0 {
		// User doesn't have a token, create a new one
		newToken := models.Token{
			UserID:    user.ID,
			Token:     t,
			ExpiresAt: tokenExpiry,
		}
		db.Create(&newToken)
	} else {
		// User already has a token, update it
		existingToken.Token = t
		existingToken.ExpiresAt = tokenExpiry
		db.Save(&existingToken)
	}

	return c.JSON(models.Loginresponse{
		Token:  t,
		Claims: claims,
	})
}

func GetUserProfileByID(c *fiber.Ctx) error {
	var user models.User
	db := database.DB.Db
	id, _ := strconv.Atoi(c.Params("id"))
	// If the user is authorized, proceed to fetch the user profile
	err := models.QueryAndReturnError(c, db, &user, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})

	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not find user", err)
	}

	var userProfile models.UserProfile
	copier.Copy(&userProfile, &user)
	return c.JSON(userProfile)
}

func UpdateUserProfileByID(c *fiber.Ctx) error {
	var user models.User
	db := database.DB.Db
	id, _ := strconv.Atoi(c.Params("id"))

	err := models.QueryAndReturnError(c, db, &user, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})

	if err != nil {
		handleError(c, fiber.StatusInternalServerError, "Could not find user", err)
	}

	// Copy contents of user to the new variable
	newUser := &user

	// Parse the request body into newUser
	if err := c.BodyParser(newUser); err != nil {
		return handleError(c, fiber.StatusBadRequest, "Invalid data", err)
	}
	db.Save(&newUser)
	var userProfile models.UserProfile
	copier.Copy(&userProfile, &newUser)
	return c.JSON(userProfile)
}

func DeleteUserByID(c *fiber.Ctx) error {
	var user models.User
	db := database.DB.Db
	id, _ := strconv.Atoi(c.Params("id"))

	err := models.QueryAndReturnError(c, db, &user, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})

	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not find user", err)
	}

	if err := db.Delete(&user).Error; err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Error deleting user", err)
	}

	return c.JSON(fiber.Map{"detail": "success"})
}

func GetUserFollowers(c *fiber.Ctx) error {
	db := database.DB.Db
	var followers []models.UserFollower

	// Query the database to find followers and preload the Source and Target users
	if err := db.Where("source_id = ?", c.Params("id")).Preload("Source").Preload("Target").Find(&followers).Error; err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not retrieve followers", err)
	}

	return c.JSON(fiber.Map{
		"followers": followers,
	})
}

func FollowUser(c *fiber.Ctx) error {
	var sourceUser, targetUser models.User
	db := database.DB.Db
	var userFollower models.UserFollower
	sourceID, _ := strconv.Atoi(c.Params("id"))
	targetID, _ := strconv.Atoi(c.Params("targetID"))

	err := db.First(&sourceUser, sourceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handleError(c, fiber.StatusNotFound, "Source user not found", err)
	}

	err = db.First(&targetUser, targetID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return handleError(c, fiber.StatusNotFound, "Source user not found", err)
	}

	userFollower = models.UserFollower{
		SourceID: uint(sourceID),
		Source:   sourceUser,
		TargetID: uint(targetID),
		Target:   targetUser,
		Type:     0, //0 - basic type of follow, for now
	}
	if err := db.Create(&userFollower).Error; err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not create follower", err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"detail":   "Followed successfully",
		"follower": userFollower,
	})
}

func UnfollowUser(c *fiber.Ctx) error {
	sourceID, _ := strconv.Atoi(c.Params("id"))
	targetID, _ := strconv.Atoi(c.Params("targetID"))

	db := database.DB.Db
	var userFollower models.UserFollower

	err := models.QueryAndReturnError(c, db, &userFollower, func(db *gorm.DB) *gorm.DB {
		return db.Where("source_id = ?", sourceID).Where("target_id = ?", targetID)
	})

	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not find record.", err)
	}

	if err := db.Delete(&userFollower).Error; err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not delete record", err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"detail": "Record deleted succesfully.",
	})
}

// The handleError function allows for customizable and quick error formatting.
func handleError(c *fiber.Ctx, statusCode int, message string, err error) error {
	if err != nil {
		return c.Status(statusCode).JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}
	return c.Status(statusCode).JSON(fiber.Map{
		"message": message,
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
