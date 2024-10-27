package transport

import (
	"api/service/workspace_user"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

// getWorkspaceUserInvitationNotVerifiedList godoc
// @Summary Get workspace user invitation not verified list (X-User-Email required, X-Workspace-Id required)
// @Description Get workspace user invitation not verified list (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Success 200 {array} workspace_user_dtos.GetWorkspaceUserListResponse
// @Router /api/v1/workspace_user/get-workspace_user_invitation_not_verified_list [get]
func (h *WorkspaceUserHandler) getWorkspaceUserInvitationNotVerifiedList(c *fiber.Ctx) error {
	workspaceUserLocal := c.Locals("workspace_user")
	if workspaceUserLocal == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	workspaceUser, ok := workspaceUserLocal.(*models.TwWorkspaceUser)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	workspaceUserStr := strconv.Itoa(workspaceUser.WorkspaceId)
	var workspaceUserList, err = workspace_user.NewWorkspaceUserService().GetWorkspaceUserInvitationNotVerifiedList(workspaceUserStr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(workspaceUserList)
}
