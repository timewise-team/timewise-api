package feature

import (
	_ "api/docs"
	"api/feature/authentication/transport"
	"api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/gofiber/swagger"
	"strings"
)

var whitelistPaths = []string{
	"/api/v1/swagger",
	"/api/v1/auth",
}

func isWhitelisted(path string) bool {
	for _, p := range whitelistPaths {
		if strings.HasPrefix(path, p) {
			println("skip middleware")
			return true
		}
	}
	return false
}

func RegisterHandlerV1() *fiber.App {
	router := fiber.New()

	// Apply AuthMiddleware with skip logic
	router.Use(skip.New(middleware.AuthMiddleware, func(ctx *fiber.Ctx) bool {
		return isWhitelisted(ctx.Path())
	}))

	// Register API v1 routes
	v1 := router.Group("/api/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)

	// Register auth routes
	transport.RegisterAuthHandler(v1.Group("/auth"))

	return router
}
