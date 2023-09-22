package main

import (
	"github.com/coaltail/GoOrders/database"
	"github.com/coaltail/GoOrders/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDb()
	app := fiber.New()

	// Set up your routes after adding the middleware
	routes.SetupUserRoutes(app)

	// Start your Fiber app
	app.Listen(":3000")
}
