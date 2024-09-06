package feature

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterHandlerV1() *fiber.App {
	router := fiber.New()
	//secretKey := authStorage.GetSecretKey()
	//router.Use(cors.Default())
	// add swagger
	v1 := router.Group("/api/v1")
	auth := v1.Group("/auth")
	{
		auth.Get("/", func(c *fiber.Ctx) error {
			return c.SendString("Hello, Fiber!")
		})
	}
	return router
}
