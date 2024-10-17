package transport

import (
	"api/service/workspace_user"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

// getWorkspaceUserInvitationList
// @Summary Get workspace user invitation list (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser (X-User-Email required, X-Workspace-Id required)
// @Produce json
// @Success 200 {array} workspace_user_dtos.GetWorkspaceUserListResponse
// @Router /api/v1/workspace_user/workspace_user_invitation_list [get]
func (h *WorkspaceUserHandler) getWorkspaceUserInvitationList(c *fiber.Ctx) error {
	workspaceUserLocal := c.Locals("workspace_user")
	if workspaceUserLocal == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	workspaceUser, ok := workspaceUserLocal.(models.TwWorkspaceUser)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	workspaceUserStr := strconv.Itoa(workspaceUser.WorkspaceId)
	workspaceUserList, err := workspace_user.NewWorkspaceUserService().GetWorkspaceUserInvitationList(workspaceUserStr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(workspaceUserList)
}
