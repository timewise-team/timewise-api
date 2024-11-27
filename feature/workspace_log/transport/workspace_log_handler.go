package transport

import (
	"api/service/workspace_log"
	"github.com/gofiber/fiber/v2"
)

// getWorkspaceLogs godoc
// @Summary Get workspace logs
// @Description Get workspace logs
// @Tags workspace_log
// @Accept json
// @Produce json
// @Param workspace_id path string true "Workspace ID"
// @Security bearerToken
// @Param X-User-Email header string true "User email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Success 200 {array} models.TwWorkspaceLog
// @Router /api/v1/workspace_log/get-workspace-logs/workspace/{workspace_id} [get]
func (h WorkspaceLogHandler) getWorkspaceLogs(ctx *fiber.Ctx) error {
	id := ctx.Params("workspace_id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "workspaceId is required",
		})
	}
	workspaceLogs, err := workspace_log.NewWorkspaceLogService().GetWorkspaceLogs(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "The server encountered an internal error",
		})
	}

	return ctx.JSON(workspaceLogs)
}
