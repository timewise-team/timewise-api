package transport

import (
	"api/service/workspace"
	workspace_utils "api/utils/workspace"
	"github.com/gofiber/fiber/v2"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/create_workspace_dtos"
)

// createWorkspace godoc
// @Summary Create a workspace
// @Description Create a workspace
// @Tags workspace
// @Accept json
// @Produce json
// @Param body body dtos.CreateWorkspaceRequest true "Create workspace request"
// @Success 201 {object} dtos.CreateWorkspaceResponse
// @Router /api/v1/workspace/create-workspace [post]
func (h *WorkspaceHandler) createWorkspace(c *fiber.Ctx) error {
	// Parse the request
	var createWorkspaceRequest dtos.CreateWorkspaceRequest
	if err := c.BodyParser(&createWorkspaceRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if err := workspace_utils.ValidateWorkspace(createWorkspaceRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Create the workspace
	workspace, err := workspace.NewCreateWorkspaceService().InitWorkspace(createWorkspaceRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if workspace == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create workspace",
		})
	}
	// Return the response
	return c.JSON(workspace)
}
