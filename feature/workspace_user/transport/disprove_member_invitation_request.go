package transport

import (
	"api/service/workspace_user"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"net/url"
)

// disproveMemberInvitationRequest godoc
// @Summary Disprove member's request invitation (X-User-Email required, X-Workspace-Id required)
// @Description Disprove member's request invitation (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Success 200 {object} models.TwWorkspaceUser
// @Router /api/v1/workspace_user/disprove-invitation/email/{email} [put]
func (h *WorkspaceUserHandler) disproveMemberInvitationRequest(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "email is required",
		})
	}
	emailFix, err1 := url.QueryUnescape(email)
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err1.Error(),
		})
	}
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
	var err = workspace_user.NewWorkspaceUserService().DisproveWorkspaceUserInvitation(workspaceUser, emailFix)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if workspaceUser == nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to verify member invitation request",
		})
	}
	return c.JSON(workspaceUser)

}
