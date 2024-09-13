package transport

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
)

func (h *AuthHandler) login(c *fiber.Ctx) error {
	user := models.TwUser{
		FirstName: "Khanh",
	}
	return c.SendString("logged in successfully for user: " + user.FirstName)
}
