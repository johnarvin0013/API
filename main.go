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
	database.DBConn.AutoMigrate(&user.User{})        //Auto Migrate - automatically make new table inside the database; Auto-Migrate came from the package gorm
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                                           //* means all user can access the API or you can just specify the user that can use the system (example: facebook.com) only
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", //header that we only accept
	}))

	// Public endpoints
	// Users
	app.Get("/api/v1/user/", user.GetUsers)           //return list of users in JSON
	app.Post("/api/v1/user/", user.NewUser)           //accepts name, username, password and role
	app.Post("/api/v1/user/auth/", user.Authenticate) //accept username and password, return token and user object in JSON

	app.Use(authRequired())

	// Authentication-required enpoints
	//API Service
	app.Get("/api/v1/ip/:ip", apiservice.GetGMapURL) //These are the 3 breaking points or changes: 1 /api/v1/ip/:ip endpoint has been changed
	// Users
	app.Put("/api/v1/user/", user.UpdateUser)       //2 you need to login first before updating
	app.Delete("/api/v1/user/:id", user.DeleteUser) //3 you need to login first before deleting

	app.Listen(":1234")
}
