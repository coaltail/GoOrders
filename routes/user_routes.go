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

	userRoutes := app.Group("/user")
	userRoutes.Post("/create", jwt, handlers.CreateUser)
	userRoutes.Post("/login", handlers.LoginUser)
	userRoutes.Get("/", jwt, handlers.ListAllUsers)
}
