package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Implement your authentication logic here
	// For example, check for a valid token in the request header
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	// Continue to the next middleware/handler
	return c.Next()
}
