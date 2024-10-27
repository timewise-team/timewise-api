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
// @Tags workspaceUser
// @Accept json
// @Produce json
// @Param workspace_user_id path int true "Workspace user ID"
// @Success 200 {array} workspace_user_dtos.GetWorkspaceUserListResponse
// @Router /dbms/v1/workspace_user/get-workspace_user_invitation_not_verified_list/workspace/{workspace_user_id} [get]
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
	if workspaceUserList == nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get workspace user list",
		})
	}
	return c.JSON(workspaceUserList)
}
