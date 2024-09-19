package feature

import (
	_ "api/docs"
	"api/feature/authentication/transport"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @host localhost:8080
// @BasePath /api/v1
func RegisterHandlerV1() *fiber.App {
	router := fiber.New()
	v1 := router.Group("/api/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)
	//v1.Use(middleware.AuthMiddleware)
	transport.RegisterAuthHandler(v1.Group("/auth"))

	return router
}
