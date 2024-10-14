package transport

import (
	"api/service/workspace_user"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

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
