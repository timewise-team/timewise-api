package transport

import (
	"api/service/board_columns"
	"github.com/gofiber/fiber/v2"
)

// getBoardColumnsByWorkspace godoc
// @Summary Get board columns by workspace (X-User-Email required, X-Workspace-Id required)
// @Description Get board columns by workspace (X-User-Email required, X-Workspace-Id required)
// @Tags board_columns
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {array} models.TwBoardColumn
// @Router /api/v1/board_columns/workspace/{workspace_id} [get]
func (h *BoardColumnsHandler) getBoardColumnsByWorkspace(c *fiber.Ctx) error {
	// Parse the request
	workspaceID := c.Params("workspace_id")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid workspace ID",
		})
	}

	// Get the board columns
	boardColumns, err := board_columns.NewBoardColumnsService().GetBoardColumnsByWorkspace(workspaceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if boardColumns == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get board columns",
		})
	}
	// Return the response
	return c.JSON(boardColumns)
}
