package transport

import (
	"api/service/board_columns"
	"github.com/gofiber/fiber/v2"
)

// deleteBoardColumn
// @Summary Delete a board column (X-User-Email required, X-Workspace-Id required)
// @Description Delete a board column (X-User-Email required, X-Workspace-Id required)
// @Tags board_columns
// @Accept json
// @Produce json
// @Param board_column_id path string true "Board column ID"
// @Param X-User-Email header string true "User email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Success 200 {object} fiber.Map
// @Router /api/v1/board_columns/{board_column_id} [delete]
func (h *BoardColumnsHandler) deleteBoardColumn(c *fiber.Ctx) error {
	// Parse the request
	boardColumnId := c.Params("board_column_id")
	if boardColumnId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Board column id is required",
		})
	}
	board_column, err1 := board_columns.NewBoardColumnsService().GetBoardColumnById(boardColumnId)
	if err1 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Delete the board column
	var err = board_columns.NewBoardColumnsService().DeleteBoardColumn(boardColumnId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	err = board_columns.NewBoardColumnsService().UpdatePositionsAfterDeletion(board_column.Position, board_column.WorkspaceId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Return the response
	return c.JSON(fiber.Map{
		"message": "Board column deleted successfully",
	})
}
