package routes

import (
	"os"

	"github.com/coaltail/GoOrders/handlers"
	"github.com/coaltail/GoOrders/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	jwt_secret := os.Getenv("JWT_SECRET")
	jwt := middlewares.NewAuthMiddleware(jwt_secret)

	userRoutes := app.Group("/users")
	userRoutes.Post("/create", jwt, handlers.CreateUser)
	app.Post("/login", handlers.LoginUser)
	userRoutes.Get("/", jwt, handlers.ListAllUsers)
	userRoutes.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return handlers.GetUserProfileByID(c, id)
	})
	userRoutes.Patch("/:id/update", func(c *fiber.Ctx) error {

		id := c.Params("id")
		return handlers.UpdateUserProfileByID(c, id)
	})
}
