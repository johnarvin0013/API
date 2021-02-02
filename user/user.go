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
	err := fiberutils.ParseBody(c, user) //(c (pagkukuhanan), user (pagpapasahan))
	username := user.Username
	password := user.Password
	db.First(user, "username = ?", user.Username) //(user , "username = ?" (?-placeholder), user.Username)

	if passwordhashing.CheckPasswordHash(password, user.Password) {
		message = fiberutils.Message{
			Message: fmt.Sprintf("User \"%s\"(%s) Authorized", user.Name, username),
			Status:  "success",
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS256) //generate a random hash token when user login, token usually used to identify or use to as a key to access some form

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["user"] = user
		claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

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
			MaxAge: 15 * 60, //(15 (minutes) * 60 (seconds))
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
	db.Create(&user) //db. is from gorm package which automatically create a query statement either create, update, delete or view //&-get reference from user
	c.JSON(fiberutils.Message{
		Message: "User Successfully Created",
		Status:  "success",
	})
	return err
}

func UpdateUser(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userO := (claims["user"].(map[string]interface{}))

	db := database.DBConn
	user := new(User)
	err := fiberutils.ParseBody(c, user)

	if userO["role"] == "admin" || fmt.Sprintf("%v", userO["id"]) == fmt.Sprintf("%d", user.Id) {
		user.Password, err = passwordhashing.HashPassword(user.Password)
		db.Model(&user).Where("id = ?", user.Id).Updates(user)
		c.JSON(fiberutils.Message{
			Message: "User Successfully Updated",
			Status:  "success",
		})
	} else {
		c.JSON(fiberutils.Message{
			Message: "No permission to update",
			Status:  "failed",
		})
	}

	return err
}

func DeleteUser(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userO := (claims["user"].(map[string]interface{}))

	db := database.DBConn
	id := c.Params("id")
	user := new(User)
	db.First(&user, id)

	if userO["role"] == "admin" || fmt.Sprintf("%v", userO["id"]) == fmt.Sprintf("%d", user.Id) {
		db.Delete(&user)
		c.JSON(fiberutils.Message{
			Message: "User Successfully Deleted",
			Status:  "success",
		})
	} else {
		c.JSON(fiberutils.Message{
			Message: "No permission to delete",
			Status:  "failed",
		})
	}
	return nil
}

func CheckToken(c *fiber.Ctx) error {
	c.SendString("Welcome to Dashboard")
	return nil
}
