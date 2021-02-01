package fiberutils

import "github.com/gofiber/fiber/v2"

type Message struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func ParseBody(c *fiber.Ctx, in interface{}) error {
	err := c.BodyParser(in)

	if err != nil {
		c.Status(503).SendString(string(err.Error()))
	}

	return err
}

func GetParamValue(c *fiber.Ctx, param string, message string) string {
	ipAddress := c.Params("ip")

	if ipAddress == "" {
		c.JSON(Message{
			Message: message,
			Status:  "failed",
		})
	}

	return ipAddress
}
