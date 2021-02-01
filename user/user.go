package user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/johnarvin0013/fiber/database"
	"github.com/johnarvin0013/fiber/fiberutils"
	"github.com/johnarvin0013/fiber/passwordhashing"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	gorm.Model
	Id       uint   `json:"id" gorm:"primarykey"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func GetUsers(c *fiber.Ctx) error {
	db := database.DBConn
	var users []User
	db.Find(&users)
	return c.JSON(users)
}

func Authenticate(c *fiber.Ctx) error {
	db := database.DBConn
	user := new(User)
	message := fiberutils.Message{}
	err := fiberutils.ParseBody(c, user)
	username := user.Username
	password := user.Password
	db.First(user, "username = ?", user.Username)

	if passwordhashing.CheckPasswordHash(password, user.Password) {
		message = fiberutils.Message{
			Message: fmt.Sprintf("User \"%s\"(%s) Authorized", user.Name, username),
			Status:  "success",
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user.Username
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		if user.Role == "admin" {
			claims["admin"] = true
		} else {
			claims["admin"] = false
		}

		t, err := token.SignedString([]byte("secret"))

		if err != nil {
			err = c.SendStatus(fiber.StatusInternalServerError)
		}

		userJSON, _ := json.Marshal(user)

		c.JSON(fiber.Map{"token": t, "user": user})
		c.Cookie(&fiber.Cookie{
			Name:   "token",
			Value:  t,
			MaxAge: 15 * 60,
		})
		c.Cookie(&fiber.Cookie{
			Name:   "user",
			Value:  string(userJSON),
			MaxAge: 15 * 60,
		})
		// return nil
	} else {
		message = fiberutils.Message{
			Message: fmt.Sprintf("User (%s) Unuthorized", username),
			Status:  "fail",
		}
		err = c.SendStatus(fiber.StatusUnauthorized)
		c.JSON(message)
	}

	// c.JSON(message)
	return err
}

func NewUser(c *fiber.Ctx) error {
	db := database.DBConn
	user := new(User)
	err := fiberutils.ParseBody(c, user)
	user.Password, err = passwordhashing.HashPassword(user.Password)
	db.Create(&user)
	c.JSON(user)
	return err
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.DBConn
	user := new(User)
	err := fiberutils.ParseBody(c, user)
	user.Password, err = passwordhashing.HashPassword(user.Password)
	db.Model(&user).Where("id = ?", user.Id).Updates(user)
	c.JSON(user)
	return err
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	user := new(User)
	db.First(&user, id)
	db.Delete(&user)
	c.JSON(user)
	return nil
}

func CheckToken(c *fiber.Ctx) error {
	c.SendString("Welcome to Dashboard")
	return nil
}
