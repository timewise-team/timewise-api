package transport

import (
	"api/service/workspace"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

// deleteWorkspace godoc
// @Summary Delete a workspace (X-User-Email required, X-Workspace-Id required)
// @Description Delete a workspace (X-User-Email required, X-Workspace-Id required)
// @Tags workspace
// @Accept json
// @Produce json
// @Security bearerToken
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Success 200 {object} fiber.Map
// @Router /api/v1/workspace/delete-workspace [delete]
func (h *WorkspaceHandler) deleteWorkspace(c *fiber.Ctx) error {
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
	workspaceIdStr := strconv.Itoa(workspaceUser.WorkspaceId)
	if workspaceIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid workspaceId"})
	}
	err := workspace.NewWorkspaceService().DeleteWorkspace(workspaceIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Delete workspace successfully"})

}
