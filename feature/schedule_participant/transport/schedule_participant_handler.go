package transport

import (
	"api/config"
	"api/service/schedule_participant"
	"api/service/workspace_user"
	auth_utils "api/utils/auth"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

type ScheduleParticipantHandler struct {
	service schedule_participant.ScheduleParticipantService
}

func NewScheduleParticipantHandler() *ScheduleParticipantHandler {
	service := schedule_participant.NewScheduleParticipantService()
	return &ScheduleParticipantHandler{
		service: *service,
	}
}

// getScheduleParticipantsByScheduleId godoc
// @Summary Get schedule participants by schedule ID
// @Description Get schedule participants by schedule ID
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param scheduleId path string true "Schedule ID"
// @Success 200 {array} schedule_participant_dtos.ScheduleParticipantInfo
// @Router /api/v1/schedule_participant/schedule/{scheduleId} [get]
func (h *ScheduleParticipantHandler) GetScheduleParticipantByScheduleID(c *fiber.Ctx) error {
	scheduleIDStr := c.Params("scheduleID")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid schedule ID")
	}
	scheduleParticipant, err := h.service.GetScheduleParticipantsByScheduleID(scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(scheduleParticipant)
}

// sendInvitation godoc
// @Summary Send invitation to user
// @Description Send invitation to user (X-User-Email required, X-Workspace-Id required)
// @Tags ScheduleParticipant
// @Accept json
// @Produce json
// @Param schedule_participant body schedule_participant_dtos.InviteToScheduleRequest true "Request body"
// @Success 200 {object} schedule_participant_dtos.ScheduleParticipantResponse
// @Router /api/v1/schedule_participant/invite [post]
func (h *ScheduleParticipantHandler) InviteToSchedule(c *fiber.Ctx) error {
	var InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest
	if err := c.BodyParser(&InviteToScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return errors.New("Failed to retrieve workspace user from context")
	}

	participant, ok := c.Locals("scheduleParticipant").(models.TwScheduleParticipant)
	if !ok {
		return fiber.NewError(500, "Failed to retrieve schedule participant")
	}
	workspaceUserInvited, err := schedule_participant.NewScheduleParticipantService().GetWorkspaceUserByEmail(
		InviteToScheduleDto.Email, workspaceUser.WorkspaceId,
	)
	if err != nil {
		return err
	}

	workspaceUser = c.Locals("workspace_user").(*models.TwWorkspaceUser)

	if workspaceUserInvited.ID != 0 {
		_, acceptLink, err1 := h.service.InviteToSchedule(workspaceUser, InviteToScheduleDto)
		if err1 != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err1.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"message":     "Invited to schedule",
			"accept_link": acceptLink,
		})
	} else {
		existingLinkedWorkspaceUser, checkErr := workspace_user.NewWorkspaceUserService().GetExistingLinkedWorkspaceUser(InviteToScheduleDto.Email, strconv.Itoa(workspaceUser.WorkspaceId))
		if checkErr != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}
		if len(existingLinkedWorkspaceUser) > 0 {
			return c.Status(500).JSON(fiber.Map{
				"message": "This email is already linked to another workspace",
				"email":   existingLinkedWorkspaceUser[0],
			})
		}

		_, _, acceptLink, err1 := h.service.InviteOutsideWorkspace(workspaceUser, participant, InviteToScheduleDto)
		if err1 != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err1.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"message":     "Invited to schedule",
			"accept_link": acceptLink,
		})
	}

	return c.JSON("invite sucessfully")
}

// assginMember godoc
// @Summary Assign member to schedule
// @Description Send invitation to user (X-User-Email required, X-Workspace-Id required)
// @Tags ScheduleParticipant
// @Accept json
// @Produce json
// @Param schedule_participant body schedule_participant_dtos.InviteToScheduleRequest true "Request body"
// @Success 200 {object} schedule_participant_dtos.ScheduleParticipantResponse
// @Router /api/v1/schedule_participant/assign [put]
func (h *ScheduleParticipantHandler) AssignMember(c *fiber.Ctx) error {
	var InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest
	if err := c.BodyParser(&InviteToScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return errors.New("Failed to retrieve workspace user from context")
	}

	workspaceUserInvited, err := schedule_participant.NewScheduleParticipantService().GetWorkspaceUserByEmail(
		InviteToScheduleDto.Email, workspaceUser.WorkspaceId,
	)
	if err != nil {
		return err
	}
	if workspaceUserInvited.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "member does not exist in workspace",
		})
	}

	scheduleParticipantInvited, err := h.service.GetScheduleParticipantsByScheduleAndWorkspaceUser(strconv.Itoa(InviteToScheduleDto.ScheduleId), strconv.Itoa(workspaceUserInvited.ID))
	if err != nil {
		return err
	}
	if scheduleParticipantInvited.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "member does not exist in schedule",
		})
	}

	workspaceUser = c.Locals("workspace_user").(*models.TwWorkspaceUser)

	scheduleParticipant, err1 := h.service.AssignMember(workspaceUser, scheduleParticipantInvited)
	if err1 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err1.Error(),
		})
	}
	return c.JSON(scheduleParticipant)
}

// acceptInvitationViaEmail godoc
// @Summary Accept invitation via email
// @Description Accept invitation via email
// @Tags ScheduleParticipant
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {object} map[string]interface{} "Schedule invitation accepted successfully"
// @Failure 404 {object} map[string]interface{} "Schedule user not found"
// @Failure 401 {object} map[string]interface{} "Token expired or invalid"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/schedule_participant/accept-invitation-via-email/token/{token} [get]
func (h *ScheduleParticipantHandler) AcceptInvite(c *fiber.Ctx) error {
	cfg, err1 := config.LoadConfig()
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to load config",
		})
	}
	token := c.Params("token")
	claims, err2 := auth_utils.ParseInvitationToken(token, cfg.JWT_SECRET)
	if err2 != nil {
	}
	scheduleId := claims["schedule_id"].(float64)
	workspaceUserId := claims["workspace_user_id"].(float64)
	scheduleIdStr := strconv.FormatFloat(scheduleId, 'f', -1, 64)
	workspaceUserIdStr := strconv.FormatFloat(workspaceUserId, 'f', -1, 64)

	// Gọi hàm AcceptInvite với các tham số đã chuyển đổi
	_, err := h.service.AcceptInvite(scheduleIdStr, workspaceUserIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Set("Content-Type", "text/html")
	return c.SendString(successHtml("already a member"))
}

// declineInvitationViaEmail godoc
// @Summary Decline invitation via email
// @Description Decline invitation via email
// @Tags ScheduleParticipant
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {object} map[string]interface{} "Invitation declined successfully"
// @Failure 404 {object} map[string]interface{} "Schedule user not found"
// @Failure 401 {object} map[string]interface{} "Invalid or expired token"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/v1/schedule_participant/decline-invitation-via-email/token/{token} [get]
func (h *ScheduleParticipantHandler) DeclineInvite(c *fiber.Ctx) error {
	cfg, err1 := config.LoadConfig()
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to load config",
		})
	}
	token := c.Params("token")
	claims, err2 := auth_utils.ParseInvitationToken(token, cfg.JWT_SECRET)
	if err2 != nil {
	}
	scheduleId := claims["schedule_id"].(float64)
	workspaceUserId := claims["workspace_user_id"].(float64)
	scheduleIdStr := strconv.FormatFloat(scheduleId, 'f', -1, 64)
	workspaceUserIdStr := strconv.FormatFloat(workspaceUserId, 'f', -1, 64)

	// Gọi hàm AcceptInvite với các tham số đã chuyển đổi
	scheduleParticipant, err := h.service.DeclineInvite(scheduleIdStr, workspaceUserIdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(scheduleParticipant)
}

// removeParticipant godoc
// @Summary Remove member from schedule
// @Description Remove participant (X-User-Email required, X-Workspace-Id required)
// @Tags ScheduleParticipant
// @Accept json
// @Produce json
// @Param schedule_participant body schedule_participant_dtos.RemoveMemberRequest true "Request body"
// @Param participant_id path string true "Participant ID"
// @Success 200 {object} schedule_participant_dtos.ScheduleParticipantResponse
// @Router /api/v1/schedule_participant/remove/{id} [put]
func (h *ScheduleParticipantHandler) RemoveParticipant(c *fiber.Ctx) error {
	participantId := c.Params("id")

	scheduleParticipant, err := h.service.RemoveParticipant(participantId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(scheduleParticipant)
}

// unassignMember godoc
// @Summary Unassign member
// @Description Unassign member (X-User-Email required, X-Workspace-Id required)
// @Tags ScheduleParticipant
// @Accept json
// @Produce json
// @Param schedule_participant body schedule_participant_dtos.RemoveMemberRequest true "Request body"
// @Param participant_id path string true "Participant ID"
// @Success 200 {object} schedule_participant_dtos.ScheduleParticipantResponse
// @Router /api/v1/schedule_participant/unassign/{id} [put]
func (h *ScheduleParticipantHandler) UnassignParticipant(c *fiber.Ctx) error {
	participantId := c.Params("id")

	scheduleParticipant, err := h.service.UnassignMember(participantId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(scheduleParticipant)
}

func errorHtml(message string) string {
	return `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Invalid Email Request</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                background-color: #f8f9fa;
                color: #343a40;
                display: flex;
                justify-content: center;
                align-items: center;
                height: 100vh;
                margin: 0;
            }
            .container {
                text-align: center;
                max-width: 500px;
                background: #fff;
                padding: 20px;
                border-radius: 10px;
                box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
            }
            h1 {
                color: #e74c3c;
                font-size: 24px;
            }
            p {
                font-size: 16px;
                line-height: 1.5;
                margin: 10px 0;
            }
            a {
                display: inline-block;
                margin-top: 20px;
                padding: 10px 20px;
                background-color: #007bff;
                color: white;
                text-decoration: none;
                border-radius: 5px;
            }
            a:hover {
                background-color: #0056b3;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Oops! Invalid Request</h1>
            <p>` + message + `</p>
            <p>If you think this is a mistake, please contact support for further assistance.</p>
        </div>
    </body>
    </html>
    `
}

func successHtml(action string) string {
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Email ` + action + ` Success</title>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { text-align: center; margin-top: 50px; }
				.success { color: green; font-size: 20px; }
				.error { color: red; font-size: 20px; }
				.button { padding: 10px 20px; font-size: 16px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1 class="success">Congratulations! Your email has been successfully ` + action + ` the workspace invitation.</h1>
				<p>If you ` + action + ` the workspace invitation, you has been added to this workspace accordingly.</p>`

	if action == "accept" {
		htmlContent += `
				<p>Your acceptance has been confirmed. You can now join the workspace.</p>`
	} else if action == "reject" {
		htmlContent += `
				<p>Your invitation has been rejected. If this was a mistake, please contact support.</p>`
	}
	htmlContent += `
				<a href="https://timewise.space/" class="button">You can close this page now.</a>
			</div>
		</body>
		</html>
	`
	return htmlContent
}
