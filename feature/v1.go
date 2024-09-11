package feature

import (
	"api/feature/authentication"
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlerV1() *fiber.App {
	router := fiber.New()
	v1 := router.Group("/api/v1")

	authentication.RegisterAuthHandler(v1.Group("/auth"))

	return router
}
