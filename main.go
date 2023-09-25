package main

import (
	"github.com/coaltail/GoOrders/database"
	"github.com/coaltail/GoOrders/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	// Set up your routes after
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"detail": "Welcome to the API!",
		})
	})
	routes.SetupUserRoutes(app)

	// Start your Fiber app
	app.Listen(":3000")
}
