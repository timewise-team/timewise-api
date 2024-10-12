package transport

import (
	"github.com/gofiber/fiber/v2"
)

type WorkspaceUserHandler struct {
	Router fiber.Router
}

func RegisterWorkspaceUserHandler(router fiber.Router) {
	workspaceUserHandler := WorkspaceUserHandler{
		Router: router,
	}

	// Register all endpoints here
	router.Get("/get-workspace_user/email/:email?/workspace_id/:workspace_id?", workspaceUserHandler.getWorkspaceUserByEmailAndWorkspace)

}
