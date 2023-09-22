package routes

import (
	"github.com/coaltail/GoOrders/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userRoutes := app.Group("/user")
	userRoutes.Post("/create", handlers.CreateUser)
	userRoutes.Get("/", handlers.ListAllUsers)
}
