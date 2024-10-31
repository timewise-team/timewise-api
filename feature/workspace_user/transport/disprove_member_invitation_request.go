package transport

import (
	"api/service/workspace_user"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"github.com/timewise-team/timewise-models/models"
)

// disproveMemberInvitationRequest godoc
// @Summary Disprove member's request invitation (X-User-Email required, X-Workspace-Id required)
// @Description Disprove member's request invitation (X-User-Email required, X-Workspace-Id required)
// @Tags WorkspaceUser
// @Accept json
// @Produce json
// @Param schedule_participant body schedule_participant_dtos.InviteToScheduleRequest true "Request body"
// @Success 200 {object} schedule_participant_dtos.ScheduleParticipantResponse
// @Router /api/v1/workspace_user/disprove-invitation [put]
func (h *WorkspaceUserHandler) disproveMemberInvitationRequest(c *fiber.Ctx) error {
	var InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest
	if err := c.BodyParser(&InviteToScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
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
	var err = workspace_user.NewWorkspaceUserService().DisproveWorkspaceUserInvitation(workspaceUser, InviteToScheduleDto.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Disprove member's request invitation successfully",
	})

}
