package main

import (
	"github.com/johnarvin0013/fiber/apiservice"
	"github.com/johnarvin0013/fiber/database"
	"github.com/johnarvin0013/fiber/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
)

func authRequired() func(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	})
}

func main() {
	database.MySQLConnect("root", "", "", "starAPI") // database.PostgreConnect("root", "", "", "starAPI")
	app := fiber.New()

	database.DBConn.AutoMigrate(&user.User{})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	//API Service
	app.Get("/api/v1/ip/:ip", apiservice.GetGMapURL)

	// Users

	app.Get("/api/v1/user/", user.GetUsers)
	app.Post("/api/v1/user/", user.NewUser)
	app.Put("/api/v1/user/", user.UpdateUser)
	app.Delete("/api/v1/user/:id", user.DeleteUser)
	app.Post("/api/v1/user/auth/", user.Authenticate)

	app.Listen(":1234")
}
