package feature

import (
	_ "api/docs"
	authTransport "api/feature/authentication/transport"
	boardColumnsTransport "api/feature/board_columns/transport"
	linkedEmailsTransport "api/feature/linked_emails/transport"
	"api/feature/schedule/transport"
	scheduleFilterTransport "api/feature/schedule_filter/transport"
	userEmailTransport "api/feature/user_email/transport"
	workspaceTransport "api/feature/workspace/transport"
	workspaceUserTransport "api/feature/workspace_user/transport"
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
			return true
		}
	}
	return false
}

func RegisterHandlerV1() *fiber.App {
	router := fiber.New()

	// Setting CORS
	router.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization,X-User-Email,X-Workspace-ID",
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
	linkedEmailsTransport.RegisterLinkedEmailsHandler(v1.Group("/user_emails"))
	transport.RegisterScheduleHandler(v1.Group("/schedules"))
	boardColumnsTransport.RegisterBoardColumnsHandler(v1.Group("/board_columns"))
	userEmailTransport.RegisterUserEmailHandler(v1.Group("/user_email"))
	workspaceUserTransport.RegisterWorkspaceUserHandler(v1.Group("/workspace_user"))
	return router
}
