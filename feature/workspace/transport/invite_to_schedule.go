package transport

import (
	"api/service/board_columns"
	"api/service/workspace"
	workspace_utils "api/utils/workspace"
	"github.com/gofiber/fiber/v2"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/create_workspace_dtos"
)

func (h *WorkspaceHandler) AdminInviteToSchedule(c *fiber.Ctx) error {
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

	err = board_columns.NewBoardColumnsService().InitBoardColumns(workspace.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Return the response
	return c.JSON(workspace)
}
