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
	//get workspace user invitation not verified list
	router.Get("/get-workspace_user_invitation_not_verified_list", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), workspaceUserHandler.getWorkspaceUserInvitationNotVerifiedList)
	//send invitation
	router.Post("/send-invitation", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), workspaceUserHandler.sendInvitation)
	//update workspace user role
	router.Put("/update-role", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), workspaceUserHandler.updateRole)
	//verify member's request invitation
	router.Put("/verify-invitation",
		middleware.CheckWorkspaceRole([]string{"owner", "admin"}),
		middleware.CheckScheduleStatus([]string{"creator"}),
		workspaceUserHandler.verifyMemberInvitationRequest)
	//disprove member's request invitation
	router.Put("/disprove-invitation", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), workspaceUserHandler.disproveMemberInvitationRequest)
	//accept invitation via email
	router.Get("/accept-invitation-via-email/token/:token?", workspaceUserHandler.acceptInvitationViaEmail)
	//decline invitation via email
	router.Get("/decline-invitation-via-email/token/:token?", workspaceUserHandler.declineInvitationViaEmail)

	router.Delete("/delete-workspace_user/workspace_user_id/:workspace_user_id?", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), workspaceUserHandler.deleteWorkspaceUser)
}
