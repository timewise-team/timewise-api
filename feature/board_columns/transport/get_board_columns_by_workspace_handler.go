package transport

import (
	"api/service/board_columns"
	"github.com/gofiber/fiber/v2"
)

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
