package feature

import (
	_ "api/docs"
	authTransport "api/feature/authentication/transport"
	scheduleFilterTransport "api/feature/schedule_filter/transport"
	workspaceTransport "api/feature/workspace/transport"
	"api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/skip"
	"github.com/gofiber/swagger"
	"strings"
)

var whitelistPaths = []string{
	"/api/v1/swagger",
	"/api/v1/auth",
	"/api/v1/workspace",
	"/api/v1/schedule",
}

func isWhitelisted(path string) bool {
	for _, p := range whitelistPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func RegisterHandlerV1() *fiber.App {
	router := fiber.New()

	// Setting CORS
	router.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Apply AuthMiddleware with skip logic
	router.Use(skip.New(middleware.AuthMiddleware, func(ctx *fiber.Ctx) bool {
		return isWhitelisted(ctx.Path())
	}))
	// Register API v1 routes
	v1 := router.Group("/api/v1")
	v1.Get("/swagger/*", swagger.HandlerDefault)

	// Register auth routes
	authTransport.RegisterAuthHandler(v1.Group("/auth"))
	workspaceTransport.RegisterWorkspaceHandler(v1.Group("/workspace"))
	scheduleFilterTransport.RegisterScheduleFilterHandler(v1.Group("/schedule"))
	return router
}
