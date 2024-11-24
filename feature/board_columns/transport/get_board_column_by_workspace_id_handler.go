package transport

import (
	"api/service/board_columns"
	"github.com/gofiber/fiber/v2"
)

// getBoardColumnsByWorkspaceId godoc
// @Summary Get board columns by workspace ID
// @Description Get board columns by workspace ID
// @Tags board_columns
// @Accept json
// @Produce json
// @Security bearerToken
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {array} models.TwBoardColumn
// @Router /api/v1/board_columns/workspace_id/{workspace_id} [get]
func (h *BoardColumnsHandler) getBoardColumnsByWorkspaceId(ctx *fiber.Ctx) error {
	workspaceId := ctx.Params("workspace_id")
	if workspaceId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid workspaceId"})
	}
	boardColumns, err := board_columns.NewBoardColumnsService().GetBoardColumnsByWorkspace(workspaceId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(boardColumns)

}
