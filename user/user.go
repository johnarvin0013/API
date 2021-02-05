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

//Start of selecting all or looking for the user available in database
func GetUsers(c *fiber.Ctx) error {
	db := database.DBConn
	var users []User
	db.Find(&users)
	return c.JSON(users)
}

//End of selecting all or looking for the user available in database

//Start of Creating New User
func NewUser(c *fiber.Ctx) error {
	db := database.DBConn
	user := new(User)
	err := fiberutils.ParseBody(c, user)
	var usernameCount int64
	db.Model(&user).Where("username = ?", user.Username).Count(&usernameCount)
	usernameExists := usernameCount > 0

	if len(user.Username) == 0 || len(user.Password) == 0 || len(user.Name) == 0 {
		fiberutils.SendJSONMessage(c, "Please Input Username, Password and Name", false, 400)
		return nil
	}

	if len(user.Username) < 3 || len(user.Password) < 8 || len(user.Name) < 3 {
		fiberutils.SendJSONMessage(c, "Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively", false, 400)
		return nil
	}

	if usernameExists {
		fiberutils.SendJSONMessage(c, "Username Already Exists", false, 400)
		return nil
	}

	user.Password, err = passwordhashing.HashPassword(user.Password)
	db.Create(&user) //db. is from gorm package which automatically create a query statement either create, update, delete or view //&-get reference from user
	fiberutils.SendJSONMessage(c, "User Successfully Created", true, 201)
	return err
}

//End of Creating New User

//Start of Login and Authenticating User
func Authenticate(c *fiber.Ctx) error {
	db := database.DBConn
	user := new(User)
	err := fiberutils.ParseBody(c, user) //(c (pagkukuhanan), user (pagpapasahan))
	username := user.Username
	password := user.Password

	if len(username) == 0 || len(password) == 0 {
		fiberutils.SendJSONMessage(c, "Please Input Username and Password", false, 400)
		return nil
	}

	db.First(user, "username = ?", user.Username) //(user , "username = ?" (?-placeholder), user.Username)

	if passwordhashing.CheckPasswordHash(password, user.Password) {
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
		fiberutils.SendJSONMessage(c, "Incorrect Username or Password", false, 401)
	}

	// c.JSON(message)
	return err
}

//End of Login and Authenticating User

//Start of Updating User
func UpdateUser(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userO := (claims["user"].(map[string]interface{}))

	db := database.DBConn
	user := new(User)
	err := fiberutils.ParseBody(c, user)

	if userO["role"] == "admin" || fmt.Sprintf("%v", userO["id"]) == fmt.Sprintf("%d", user.Id) {
		if len(user.Username) == 0 || len(user.Password) == 0 || len(user.Name) == 0 {
			fiberutils.SendJSONMessage(c, "Please Input Username, Password and Name", false, 400)
			return nil
		}

		if len(user.Username) < 3 || len(user.Password) < 8 || len(user.Name) < 3 {
			fiberutils.SendJSONMessage(c, "Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively", false, 400)
			return nil
		}

		user.Password, err = passwordhashing.HashPassword(user.Password)
		var existingUser User
		db.First(&existingUser, user.Id)

		if len(existingUser.Role) == 0 {
			fiberutils.SendJSONMessage(c, "No User exists", false, 404)
			return nil
		}

		db.Updates(&user)
		fiberutils.SendJSONMessage(c, "User Successfully Updated", true, 200)
	} else {
		fiberutils.SendJSONMessage(c, "No permission to update", false, 401)
	}

	return err
}

//End of Updating User

//Start of Deleting User
func DeleteUser(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userO := (claims["user"].(map[string]interface{}))

	db := database.DBConn
	id := c.Params("id")
	user := new(User)
	db.First(&user, id)

	if len(user.Role) == 0 {
		fiberutils.SendJSONMessage(c, "No User exists", false, 404)
		return nil
	}

	if userO["role"] == "admin" || fmt.Sprintf("%v", userO["id"]) == fmt.Sprintf("%d", user.Id) {
		db.Delete(&user)
		fiberutils.SendJSONMessage(c, "User Successfully Deleted", true, 200)
	} else {
		fiberutils.SendJSONMessage(c, "No permission to delete", false, 401)
	}
	return nil
}

//End of Deleting User
