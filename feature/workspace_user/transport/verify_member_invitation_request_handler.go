package transport

import (
	"api/config"
	"api/service/auth"
	"api/service/schedule"
	"api/service/user_email"
	"api/service/workspace_user"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

// verifyMemberInvitationRequest godoc
// @Summary Verify member's request invitation (X-User-Email required, X-Workspace-Id required)
// @Description Verify member's request invitation (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Param schedule_participant body schedule_participant_dtos.InviteToScheduleRequest true "Request body"
// @Success 200 {object} schedule_participant_dtos.ScheduleParticipantResponse
// @Router /api/v1/workspace_user/verify-invitation [put]
func (h *WorkspaceUserHandler) verifyMemberInvitationRequest(c *fiber.Ctx) error {
	var InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest
	if err := c.BodyParser(&InviteToScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	cfg, err1 := config.LoadConfig()
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to load config",
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
	userEmail, errs := user_email.NewUserEmailService().GetUserEmail(InviteToScheduleDto.Email)
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
	var err = workspace_user.NewWorkspaceUserService().VerifyWorkspaceUserInvitation(workspaceUser, InviteToScheduleDto.Email)
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
			scheduleIdStr := strconv.Itoa(InviteToScheduleDto.ScheduleId)
			scheduleInfo, _ := schedule.NewScheduleService().GetScheduleByID(scheduleIdStr)
			acceptLink, declineLink, _ := auth.GenerateInviteScheduleLinks(cfg, InviteToScheduleDto.ScheduleId, workspaceUserCheck.ID)
			content := auth.BuildScheduleInvitationContent(scheduleInfo, acceptLink, declineLink)
			subject := fmt.Sprintf("Invitation to join workspace: %s", scheduleInfo.Title)
			if err := auth.SendInvitationEmail(cfg, userEmail.Email, content, subject); err != nil {
				return c.Status(500).JSON(fiber.Map{"message": "Failed to send invitation email"})
			}
			return c.Status(200).JSON(fiber.Map{
				"message": "Invitation sent successfully",
			})
		}
	}
	return c.JSON(fiber.Map{
		"workspaceUser": workspaceUserCheck,
	})

}
