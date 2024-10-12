package transport

import (
	"api/service/workspace_user"
	auth_utils "api/utils/auth"
	"github.com/gofiber/fiber/v2"
)

// getWorkspaceUserByEmailAndWorkspace
// @Summary Get workspace user by email and workspace
// @Tags WorkspaceUser
// @Produce json
// @Param email path string true "Email"
// @Param workspace_id path string true "Workspace ID"
// @Success 200 {object} models.TwWorkspaceUser
// @Router /api/v1/workspace_user/get-workspace_user/email/{email}/workspace_id/{workspace_id} [get]
func (h *WorkspaceUserHandler) getWorkspaceUserByEmailAndWorkspace(c *fiber.Ctx) error {
	var email = c.Params("email")
	if email == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "email is required",
		})
	}
	if !auth_utils.IsValidEmail(email) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid email",
		})
	}
	var workspaceId = c.Params("workspace_id")
	if workspaceId == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "workspace_id is required",
		})
	}
	workspaceUser, err := workspace_user.NewWorkspaceUserService().GetWorkspaceUserByEmailAndWorkspaceID(email, workspaceId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if workspaceUser == nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "workspace user not found",
		})
	}
	return c.JSON(workspaceUser)

}
