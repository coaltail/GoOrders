package handlers

import (
	"strconv"

	"github.com/coaltail/GoOrders/database"
	"github.com/coaltail/GoOrders/middlewares"
	"github.com/coaltail/GoOrders/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllMessages(c *fiber.Ctx) error {
	db := database.DB.Db
	var messages []models.Message
	claims, valid := middlewares.ExtractClaims(c.Get("Authorization"))
	if !valid {
		return handleError(c, fiber.StatusInternalServerError, "Could not parse claims", fiber.ErrInternalServerError)
	}
	user, err := middlewares.GetUserFromClaims(*claims, db)
	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not get user from claims.", err)
	}

	if err := models.QueryAndReturnError(c, db, &messages, func(db *gorm.DB) *gorm.DB {
		return db.Where("message_sender_id = ?", user.ID)
	}); err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not fetch messages", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": messages,
	})
}

func GetUserChat(c *fiber.Ctx) error {
	messageRecipientID, _ := strconv.Atoi(c.Params("recipientID"))
	db := database.DB.Db
	var messages []models.Message

	claims, valid := middlewares.ExtractClaims(c.Get("Authorization"))
	if !valid {
		return handleError(c, fiber.StatusInternalServerError, "Could not parse claims", fiber.ErrInternalServerError)
	}
	user, err := middlewares.GetUserFromClaims(*claims, db)
	if err != nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not get user from claims.", err)
	}

	if err := db.Where("message_sender_id = ?", user.ID).Where("message_recipient_id = ?", messageRecipientID).Find(&messages).Error; err!= nil {
		return handleError(c, fiber.StatusInternalServerError, "Could not retrieve messages", fiber.ErrInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"messages": messages,
	})

}
