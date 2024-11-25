package transport

import (
	"api/middleware"
	"github.com/gofiber/fiber/v2"
)

type WorkspaceLogHandler struct {
	Router fiber.Router
}

func RegisterWorkspaceLogeHandler(router fiber.Router) {
	workspaceLogHandler := WorkspaceLogHandler{
		Router: router,
	}
	router.Get("/get-workspace-logs/workspace/:workspace_id", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member", "guest"}), workspaceLogHandler.getWorkspaceLogs)

}
