package transport

import (
	"api/service/workspace_user"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/workspace_user_dtos"
	"github.com/timewise-team/timewise-models/models"
)

// updateRole
// @Summary Update role of workspace user (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Param workspace_user_id path string true "workspace_user_id"
// @Param workspace_user body workspace_user_dtos.UpdateWorkspaceUserRoleRequest true "Update role request"
// @Success 200 {object} fiber.Map
// @Router /api/v1/workspace_user/update-role [put]
func (s *WorkspaceUserHandler) updateRole(c *fiber.Ctx) error {
	workspaceUserLocal := c.Locals("workspace_user")
	if workspaceUserLocal == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	workspaceUser, ok := workspaceUserLocal.(*models.TwWorkspaceUser)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	var UpdateWorkspaceUserRoleRequest workspace_user_dtos.UpdateWorkspaceUserRoleRequest
	if err := c.BodyParser(&UpdateWorkspaceUserRoleRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}
	email := c.Locals("email")
	if email == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid session",
		})
	}
	err := workspace_user.NewWorkspaceUserService().UpdateWorkspaceUserRole(workspaceUser, UpdateWorkspaceUserRoleRequest, email.(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Update role successfully",
	})
}
