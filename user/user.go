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
	Role     string `json:"role" gorm:"default:user"`
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

	if len(username) == 0 || len(password) == 0 {
		c.JSON(fiberutils.Message{
			Message: "Please Input Username and Password",
			Status:  "error",
		})
		return nil
	}

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
	var usernameCount int64
	db.Model(&user).Where("username = ?", user.Username).Count(&usernameCount)
	usernameExists := usernameCount > 0

	if len(user.Username) == 0 || len(user.Password) == 0 || len(user.Name) == 0 {
		c.JSON(fiberutils.Message{
			Message: "Please Input Username, Password and Name",
			Status:  "error",
		})
		return nil
	}

	if len(user.Username) < 3 || len(user.Password) < 8 || len(user.Name) < 3 {
		c.JSON(fiberutils.Message{
			Message: "Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively",
			Status:  "error",
		})
		return nil
	}

	if usernameExists {
		c.JSON(fiberutils.Message{
			Message: "Username Already Exists",
			Status:  "error",
		})
		return nil
	}

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
		if len(user.Username) == 0 || len(user.Password) == 0 || len(user.Name) == 0 {
			c.JSON(fiberutils.Message{
				Message: "Please Input Username, Password and Name",
				Status:  "error",
			})
			return nil
		}

		if len(user.Username) < 3 || len(user.Password) < 8 || len(user.Name) < 3 {
			c.JSON(fiberutils.Message{
				Message: "Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively",
				Status:  "error",
			})
			return nil
		}

		user.Password, err = passwordhashing.HashPassword(user.Password)
		db.First(&user, user.Id)

		if len(user.Role) == 0 {
			c.JSON(fiberutils.Message{
				Message: fmt.Sprintf("No User with Id: %d", user.Id),
				Status:  "error",
			})
			return nil
		}

		db.Updates(&user)
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

	if len(user.Role) == 0 {
		c.JSON(fiberutils.Message{
			Message: "No User with Id: " + id,
			Status:  "error",
		})
		return nil
	}

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
