package transport

import (
	"api/middleware"
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
	//get workspace user list
	router.Get("/workspace_user_list", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}), workspaceUserHandler.getWorkspaceUserList)
	//get workspace user invitation list
	router.Get("/workspace_user_invitation_list", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), workspaceUserHandler.getWorkspaceUserInvitationList)

}
