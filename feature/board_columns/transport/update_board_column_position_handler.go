package transport

import (
	"api/service/board_columns"
	"github.com/gofiber/fiber/v2"
)

type updatePositionRequest struct {
	Position int `json:"position"`
}

// updatePosition godoc
// @Summary Update a board column position (X-User-Email required, X-Workspace-Id required)
// @Description Update a board column position (X-User-Email required, X-Workspace-Id required)
// @Tags board_columns
// @Accept json
// @Produce json
// @Param board_column_id path string true "Board column ID"
// @Param body body updatePositionRequest true "Update board column position request"
// @Param X-User-Email header string true "User email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Router /api/v1/board_columns/update_position/{board_column_id} [put]
func (h *BoardColumnsHandler) updatePosition(c *fiber.Ctx) error {
	// Parse the request
	boardColumnId := c.Params("board_column_id")
	if boardColumnId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Board column id is required",
		})
	}
	var boardColumn updatePositionRequest
	if err := c.BodyParser(&boardColumn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	if boardColumn.Position <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid position",
		})
	}
	oldBoardColumn, err := board_columns.NewBoardColumnsService().GetBoardColumnById(boardColumnId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	// Update the board column
	err1 := board_columns.NewBoardColumnsService().UpdatePositionAfterDrag(oldBoardColumn.Position, boardColumn.Position, oldBoardColumn.WorkspaceId, boardColumnId)
	if err1 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Return the response
	return c.JSON(fiber.Map{
		"message": "Board column position updated successfully",
	})
}
