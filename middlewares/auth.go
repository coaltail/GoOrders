package middlewares

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/coaltail/GoOrders/models"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	})
}

func CompareJWTandUserIDMiddleware() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Get the user ID from the route parameters or wherever it is available
		id, _ := strconv.Atoi(c.Params("id"))

		// Call ParseAndCompare to check if the user is authorized
		err := ParseAndCompare(id, c)
		if err != nil {
			return err // If not authorized, return the error response
		}

		// If authorized, continue to the next handler
		return c.Next()
	}
}

func ParseAndCompare(userID int, c *fiber.Ctx) error {
	claims, valid := ExtractClaims(c.Get("Authorization"))
	if !valid {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not parse token")
	}
	userIDFromJWT := claims.ID

	// Data types have to be matching, convert both source and jwt id to int
	jwtUserID := int(userIDFromJWT)
	if jwtUserID != userID {
		return fiber.NewError(fiber.StatusForbidden, "You are not authorized to make this request.")
	}

	return nil
}

// The ExtractClaims function takes in a JWT token string and extracts its claims. The claims included are: ID, user email and token expiry time.
func ExtractClaims(tokenStr string) (*models.Claims, bool) {
	// Remove the "Bearer " prefix if it exists
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	var hmacSecret = []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}
func GetUserFromClaims(claims models.Claims, db *gorm.DB) (*models.User, error) {
	//Get claims
	userID := claims.ID
	var user models.User
	//Get user from database
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	// No error, return user
	return &user, nil
}
