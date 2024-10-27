package transport

import (
	"api/config"
	"api/service/auth"
	"api/service/user_email"
	"api/service/workspace"
	"api/service/workspace_user"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"net/url"
	"strconv"
)

// verifyMemberInvitationRequest godoc
// @Summary Verify member's request invitation (X-User-Email required, X-Workspace-Id required)
// @Description Verify member's request invitation (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Success 200 {object} models.TwWorkspaceUser
// @Router /api/v1/workspace_user/verify-invitation/email/{email} [put]
func (h *WorkspaceUserHandler) verifyMemberInvitationRequest(c *fiber.Ctx) error {
	cfg, err1 := config.LoadConfig()

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
	userEmail, errs := user_email.NewUserEmailService().GetUserEmail(emailFix)
	if userEmail == nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "This email is not registered",
		})
	}
	if errs != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	var err = workspace_user.NewWorkspaceUserService().VerifyWorkspaceUserInvitation(workspaceUser, emailFix)
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
	workspaceUserCheck, err1 := workspace_user.NewWorkspaceUserService().GetWorkspaceUserByEmailAndWorkspaceID(userEmail.Email, strconv.Itoa(workspaceUser.WorkspaceId))
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	if workspaceUserCheck != nil {
		if workspaceUserCheck.Status == "pending" {
			worspaceIdStr := strconv.Itoa(workspaceUser.WorkspaceId)
			workspaceInfo := workspace.NewWorkspaceService().GetWorkspaceById(worspaceIdStr)
			acceptLink, declineLink, _ := auth.GenerateInviteLinks(cfg, userEmail.Email, workspaceInfo.ID, workspaceUserCheck.Role)
			content := auth.BuildInvitationContent(workspaceInfo, workspaceUserCheck.Role, acceptLink, declineLink)
			subject := fmt.Sprintf("Invitation to join workspace: %s", workspaceInfo.Title)
			if err := auth.SendInvitationEmail(cfg, userEmail.Email, content, subject); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
			}
			return c.Status(200).JSON(fiber.Map{
				"message": "Invitation sent successfully",
			})
		}

	}

	return c.JSON(workspaceUserCheck)

}
