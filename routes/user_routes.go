package routes

import (
	"os"

	"github.com/coaltail/GoOrders/handlers"
	"github.com/coaltail/GoOrders/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	protect_Route_secret := os.Getenv("JWT_SECRET")
	protect_Route := middlewares.NewAuthMiddleware(protect_Route_secret)

	userRoutes := app.Group("/users")
	userRoutes.Post("/create", handlers.CreateUser)
	app.Post("/login", handlers.LoginUser)
	userRoutes.Get("/", protect_Route, handlers.ListAllUsers)
	userRoutes.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return handlers.GetUserProfileByID(c, id)
	})
	userRoutes.Patch("/:id/update", func(c *fiber.Ctx) error {

		id := c.Params("id")
		return handlers.UpdateUserProfileByID(c, id)
	})
	userRoutes.Delete("/:id/delete", protect_Route, func(c *fiber.Ctx) error {
		id := c.Params("id")
		return handlers.DeleteUserByID(c, id)
	})
	userRoutes.Get("/:id/followers", protect_Route, handlers.GetUserFollowers)

}
