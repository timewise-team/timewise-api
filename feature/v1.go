package feature

import (
	_ "api/docs"
	accountTransport "api/feature/account/transport"
	authTransport "api/feature/authentication/transport"
	boardColumnsTransport "api/feature/board_columns/transport"
	commentsTransport "api/feature/comment/transport"
	documentTransport "api/feature/document/transport"
	"api/feature/schedule/transport"
	scheduleFilterTransport "api/feature/schedule_filter/transport"
	scheduleLog "api/feature/schedule_log/transport"
	scheduleParticipant "api/feature/schedule_participant/transport"
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
	"/api/v1/workspace_user/accept-invitation-via-email",
	"/api/v1/workspace_user/decline-invitation-via-email",
	"/api/v1/account/user/emails/link",
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
	transport.RegisterScheduleHandler(v1.Group("/schedules"))
	boardColumnsTransport.RegisterBoardColumnsHandler(v1.Group("/board_columns"))
	userEmailTransport.RegisterUserEmailHandler(v1.Group("/user_email"))
	workspaceUserTransport.RegisterWorkspaceUserHandler(v1.Group("/workspace_user"))
	documentTransport.RegisterDocumentHandler(v1.Group("/document"))
	scheduleParticipant.RegisterScheduleParticipantHandler(v1.Group("/schedule_participant"))
	scheduleLog.RegisterScheduleLogHandler(v1.Group("/schedule_log"))
	commentsTransport.RegisterCommentHandler(v1.Group("/comment"))
	accountTransport.RegisterAccountHandler(v1.Group("/account"))

	return router
}
