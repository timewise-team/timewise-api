package transport

import (
	"api/service/workspace_user"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

// getWorkspaceUserList
// @Summary Get workspace user list (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Produce json
// @Success 200 {array} workspace_user_dtos.GetWorkspaceUserListResponse
// @Router /api/v1/workspace_user/workspace_user_list [get]
func (h *WorkspaceUserHandler) getWorkspaceUserList(c *fiber.Ctx) error {

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
	workspaceUserList, err := workspace_user.NewWorkspaceUserService().GetWorkspaceUserList(workspaceUserStr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(workspaceUserList)
}

// getWorkspaceUserListForManage godoc
// @Summary Get workspace user list for manage (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Produce json
// @Success 200 {array} workspace_user_dtos.GetWorkspaceUserListResponse
// @Router /api/v1/workspace_user/manage/wsp_user_list [get]
func (h *WorkspaceUserHandler) getWorkspaceUserListForManage(c *fiber.Ctx) error {

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
	workspaceUserList, err := workspace_user.NewWorkspaceUserService().GetWorkspaceUserListForManage(workspaceUserStr)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(workspaceUserList)
}
