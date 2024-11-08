package transport

import (
	"api/service/board_columns"
	board_columns_utils "api/utils/board_columns"
	"github.com/gofiber/fiber/v2"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
)

// createBoardColumn godoc
// @Summary Create a board column (X-User-Email required, X-Workspace-Id required)
// @Description Create a board column (X-User-Email required, X-Workspace-Id required)
// @Tags board_columns
// @Accept json
// @Produce json
// @Param body body dtos.BoardColumnsRequest true "Create board column request"
// @Param X-User-Email header string true "User email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Success 201 {object} models.TwBoardColumn
// @Router /api/v1/board_columns [post]
func (h *BoardColumnsHandler) createBoardColumn(c *fiber.Ctx) error {
	// Parse the request
	var createBoardColumnRequest dtos.BoardColumnsRequest
	if err := c.BodyParser(&createBoardColumnRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := board_columns_utils.ValidateBoardColumn(createBoardColumnRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Create the board cgolumn
	var boardColumn, err = board_columns.NewBoardColumnsService().CreateBoardColumn(createBoardColumnRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if boardColumn == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create board column",
		})
	}
	// Return the response
	return c.JSON(boardColumn)
}
