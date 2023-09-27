package routes

import (
	"os"

	"github.com/coaltail/GoOrders/handlers"
	"github.com/coaltail/GoOrders/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	app.Post("/login", handlers.LoginUser)

	protect_Route_secret := os.Getenv("JWT_SECRET")
	protect_Route := middlewares.NewAuthMiddleware(protect_Route_secret)
	userRoutes := app.Group("/users")
	userRoutes.Post("/create", handlers.CreateUser)
	userRoutes.Get("/", protect_Route, handlers.ListAllUsers)
	userRoutes.Get("/:id", protect_Route, middlewares.CompareJWTandUserIDMiddleware(), handlers.GetUserProfileByID)
	userRoutes.Patch("/:id/update", protect_Route, middlewares.CompareJWTandUserIDMiddleware(), handlers.UpdateUserProfileByID)
	userRoutes.Delete("/:id/delete", protect_Route, middlewares.CompareJWTandUserIDMiddleware(), handlers.DeleteUserByID)

	userRoutes.Get("/:id/followers", protect_Route, handlers.GetUserFollowers)
	userRoutes.Post("/:id/followers/:targetID", protect_Route, handlers.FollowUser)
	userRoutes.Delete("/:id/followers/:targetID", protect_Route, handlers.UnfollowUser)

	userRoutes.Get("/:id/friends", protect_Route, handlers.GetUserFriends)
	userRoutes.Post("/:id/friends/:targetID", protect_Route, handlers.CreateUserFriends)
	userRoutes.Delete("/:id/friends/:targetID", protect_Route, handlers.DeleteUserFriends)

}
