package authentication

import (
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) login(c *fiber.Ctx) error {
	return c.SendString("login")
}
