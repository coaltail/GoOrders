package main

import (
	"github.com/coaltail/GoOrders/database"
	"github.com/coaltail/GoOrders/routes"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"os"
)

func main() {
	database.ConnectDb()
	app := fiber.New()
	jwt_secret := os.Getenv("JWT_SECRET")
    app.Use(jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{Key: []byte(jwt_secret)},
    }))

	routes.SetupUserRoutes(app)
	app.Listen(":3000")
}
