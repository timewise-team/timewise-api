package feature

import (
	_ "api/docs"
	"api/feature/authentication/transport"
	schedule_filter_transport "api/feature/schedule_filter/transport"
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

	// Register service routes
	scheduleFilterService := schedule_filter_transport.NewScheduleFilterRequest()

	// Register auth routes
	transport.RegisterAuthHandler(v1.Group("/auth"))
	schedule_filter_transport.RegisterScheduleFilterHandler(v1.Group("/schedule"), scheduleFilterService)
	return router
}
