package routes

import (
	"os"

	"github.com/coaltail/GoOrders/handlers"
	"github.com/coaltail/GoOrders/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupMessageRoutes(app *fiber.App) {
	protect_Route_secret := os.Getenv("JWT_SECRET")
	protect_Route := middlewares.NewAuthMiddleware(protect_Route_secret)
	messageRotues := app.Group("/messages")
	messageRotues.Get("/", protect_Route, handlers.GetAllMessages)
}
