package transport

import (
	"api/service/workspace"
	"github.com/gofiber/fiber/v2"
)

// getWorkspaceById godoc
// @Summary Get a workspace by ID
// @Description Get a workspace by ID
// @Tags workspace
// @Accept json
// @Produce json
// @Security bearerToken
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {object} models.TwWorkspace
// @Router /api/v1/workspace/get-workspace-by-id/{workspace_id} [get]
func (h *WorkspaceHandler) getWorkspaceById(c *fiber.Ctx) error {
	workspaceId := c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid workspaceId"})
	}
	workspace := workspace.NewWorkspaceService().GetWorkspaceById(workspaceId)
	if workspace == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Workspace not found"})
	}
	return c.Status(fiber.StatusOK).JSON(workspace)
}
