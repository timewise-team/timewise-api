package transport

import (
	"api/service/board_columns"
	board_columns_utils "api/utils/board_columns"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
)

// updateBoardColumn godoc
// @Summary Update a board column (X-User-Email required, X-Workspace-Id required)
// @Description Update a board column (X-User-Email required, X-Workspace-Id required)
// @Tags board_columns
// @Accept json
// @Produce json
// @Param board_column_id path string true "Board column ID"
// @Param body body board_columns_dtos.BoardColumnsRequest true "Update board column request"
// @Router /api/v1/board_columns/{board_column_id} [put]
func (h *BoardColumnsHandler) updateBoardColumn(c *fiber.Ctx) error {
	// Parse the request
	boardColumnId := c.Params("board_column_id")
	if boardColumnId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Board column id is required",
		})
	}

	// Parse the request body
	var boardColumnData board_columns_dtos.BoardColumnsRequest
	if err := c.BodyParser(&boardColumnData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
		})
	}
	if err := board_columns_utils.ValidateBoardColumn(boardColumnData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// Update the board column
	boardColumnResponse, err := board_columns.NewBoardColumnsService().UpdateBoardColumn(boardColumnId, boardColumnData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if boardColumnResponse == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update board column",
		})
	}
	if boardColumnResponse.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Board column not found"})
	}
	// Return the response
	return c.JSON(boardColumnResponse)

}
