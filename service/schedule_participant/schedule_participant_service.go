package schedule_participant

import (
	"api/config"
	"api/dms"
	"api/service/auth"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"github.com/timewise-team/timewise-models/models"
	"net/http"
	"strconv"
	"time"
)

type ScheduleParticipantService struct {
}

func NewScheduleParticipantService() *ScheduleParticipantService {
	return &ScheduleParticipantService{}
}

func (h *ScheduleParticipantService) GetScheduleParticipantsBySchedule(scheduleId int, workspaceId string) ([]schedule_participant_dtos.ScheduleParticipantInfo, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	if workspaceId == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/workspace/"+workspaceId+"/schedule/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleParticipants []schedule_participant_dtos.ScheduleParticipantInfo
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	return scheduleParticipants, nil

}

func (h *ScheduleParticipantService) GetScheduleParticipantsByScheduleID(scheduleId int) ([]schedule_participant_dtos.ScheduleParticipantInfo, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/schedule/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleParticipants []schedule_participant_dtos.ScheduleParticipantInfo
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	return scheduleParticipants, nil

}

func (h *ScheduleParticipantService) InviteToSchedule(
	c *fiber.Ctx,
	InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest,
) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	workspaceUserInvite, err := h.getWorkspaceUserFromContext(c)
	if err != nil {
		return nil, err
	}

	workspaceUserInvited, err := h.getWorkspaceUserByEmail(
		InviteToScheduleDto.Email, workspaceUserInvite.WorkspaceId,
	)
	if err != nil {
		return nil, err
	}

	scheduleParticipant, err := h.getScheduleParticipantFromContext(c)
	if err != nil {
		return nil, err
	}

	schedule, err := h.getScheduleById(scheduleParticipant.ScheduleId)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fiber.NewError(500, "Failed to load config")
	}

	acceptLink, declineLink, _ := auth.GenerateInviteScheduleLinks(
		cfg, scheduleParticipant.ScheduleId, workspaceUserInvited.ID,
	)

	scheduleParticipantResponse, err := h.handleInvitation(
		cfg, schedule, scheduleParticipant, workspaceUserInvite, workspaceUserInvited,
		acceptLink, declineLink, InviteToScheduleDto.Email,
	)

	if err != nil {
		return nil, err
	}

	return scheduleParticipantResponse, nil
}

// Hàm hỗ trợ lấy Workspace User từ context
func (h *ScheduleParticipantService) getWorkspaceUserFromContext(c *fiber.Ctx) (*models.TwWorkspaceUser, error) {
	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return nil, errors.New("Failed to retrieve workspace user from context")
	}
	return workspaceUser, nil
}

// Gọi API để lấy Workspace User theo email
func (h *ScheduleParticipantService) getWorkspaceUserByEmail(email string, workspaceId int) (*models.TwWorkspaceUser, error) {
	resp, err := dms.CallAPI(
		"GET", fmt.Sprintf("/workspace_user/email/%s/workspace/%d", email, workspaceId),
		nil, nil, nil, 120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user models.TwWorkspaceUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Lấy Schedule Participant từ context
func (h *ScheduleParticipantService) getScheduleParticipantFromContext(c *fiber.Ctx) (models.TwScheduleParticipant, error) {
	participant, ok := c.Locals("scheduleParticipant").(models.TwScheduleParticipant)
	if !ok {
		return models.TwScheduleParticipant{}, fiber.NewError(500, "Failed to retrieve schedule participant")
	}
	return participant, nil
}

// Lấy thông tin Schedule qua API
func (h *ScheduleParticipantService) getScheduleById(scheduleId int) (*models.TwSchedule, error) {
	resp, err := dms.CallAPI(
		"GET", fmt.Sprintf("/schedule/%d", scheduleId),
		nil, nil, nil, 120,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var schedule models.TwSchedule
	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		return nil, fmt.Errorf("failed to decode schedule: %v", err)
	}
	return &schedule, nil
}

// Xử lý lời mời và gửi email
func (h *ScheduleParticipantService) handleInvitation(
	cfg *config.Config, schedule *models.TwSchedule, scheduleParticipant models.TwScheduleParticipant,
	workspaceUserInvite *models.TwWorkspaceUser, workspaceUserInvited *models.TwWorkspaceUser,
	acceptLink, declineLink, email string,
) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	resp, err := dms.CallAPI(
		"GET", fmt.Sprintf("/schedule_participant/workspace_user/%d/schedule/%d",
			workspaceUserInvited.ID, scheduleParticipant.ScheduleId),
		nil, nil, nil, 120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	now := time.Now()
	if resp.StatusCode == http.StatusNotFound {
		return h.createAndSendInvitation(
			cfg, schedule, scheduleParticipant, workspaceUserInvite, workspaceUserInvited,
			acceptLink, declineLink, email, now,
		)
	}

	var existingParticipant models.TwScheduleParticipant
	if err := json.NewDecoder(resp.Body).Decode(&existingParticipant); err != nil {
		return nil, err
	}

	return h.handleExistingInvitation(
		cfg, schedule, &existingParticipant, acceptLink, declineLink, email,
	)
}

// Tạo lời mời mới và gửi email
func (h *ScheduleParticipantService) createAndSendInvitation(
	cfg *config.Config, schedule *models.TwSchedule, scheduleParticipant models.TwScheduleParticipant,
	workspaceUserInvite *models.TwWorkspaceUser, workspaceUserInvited *models.TwWorkspaceUser,
	acceptLink, declineLink, email string, now time.Time,
) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	newParticipant := models.TwScheduleParticipant{
		CreatedAt:        now,
		UpdatedAt:        now,
		ScheduleId:       scheduleParticipant.ScheduleId,
		WorkspaceUserId:  workspaceUserInvited.ID,
		AssignBy:         workspaceUserInvite.ID,
		InvitationSentAt: &now,
		InvitationStatus: "pending",
	}

	resp, err := dms.CallAPI("POST", "/schedule_participant", newParticipant, nil, nil, 120)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response schedule_participant_dtos.ScheduleParticipantResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	subject := fmt.Sprintf("Invitation to join schedule: %s", schedule.Title)
	content := auth.BuildScheduleInvitationContent(schedule, acceptLink, declineLink)
	if err := auth.SendInvitationEmail(cfg, email, content, subject); err != nil {
		return nil, fiber.NewError(500, "Failed to send invitation email")
	}

	return &response, nil
}

// Xử lý lời mời đã tồn tại
func (h *ScheduleParticipantService) handleExistingInvitation(
	cfg *config.Config, schedule *models.TwSchedule, participant *models.TwScheduleParticipant,
	acceptLink, declineLink, email string,
) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	if participant.InvitationStatus == "joined" {
		return nil, errors.New("User is already in the schedule")
	}

	subject := fmt.Sprintf("Reminder: Invitation to join schedule: %s", schedule.Title)
	content := auth.BuildScheduleInvitationContent(schedule, acceptLink, declineLink)
	if err := auth.SendInvitationEmail(cfg, email, content, subject); err != nil {
		return nil, fiber.NewError(500, "Failed to send reminder email")
	}

	if participant.InvitationStatus == "declined" || participant.InvitationStatus == "removed" {
		participant.InvitationStatus = "pending"
		participant.UpdatedAt = time.Now()

		resp, err := dms.CallAPI(
			"PUT", fmt.Sprintf("/schedule_participant/%d", participant.ID),
			participant, nil, nil, 120,
		)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	return &schedule_participant_dtos.ScheduleParticipantResponse{}, nil
}

func (h *ScheduleParticipantService) AcceptInvite(scheduleId, workspaceUserId string) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/workspace_user/"+workspaceUserId+"/schedule/"+scheduleId,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleParticipants models.TwScheduleParticipant
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	now := time.Now()
	scheduleParticipants.Status = "participant"
	scheduleParticipants.ResponseTime = &now
	scheduleParticipants.InvitationStatus = "joined"

	resp1, err := dms.CallAPI(
		"PUT",
		"/schedule_participant/"+strconv.Itoa(scheduleParticipants.ID),
		scheduleParticipants,
		nil,
		nil,
		120,
	)

	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	var updateScheduleParticipants schedule_participant_dtos.ScheduleParticipantResponse
	if err := json.NewDecoder(resp1.Body).Decode(&updateScheduleParticipants); err != nil {
		return nil, err
	}

	return &updateScheduleParticipants, nil

}

func (h *ScheduleParticipantService) DeclineInvite(scheduleId, workspaceUserId string) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/workspace_user/"+workspaceUserId+"/schedule/"+scheduleId,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleParticipants models.TwScheduleParticipant
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	now := time.Now()
	scheduleParticipants.ResponseTime = &now
	scheduleParticipants.InvitationStatus = "declined"

	resp1, err := dms.CallAPI(
		"PUT",
		"/schedule_participant/"+strconv.Itoa(scheduleParticipants.ID),
		scheduleParticipants,
		nil,
		nil,
		120,
	)

	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	var updateScheduleParticipants schedule_participant_dtos.ScheduleParticipantResponse
	if err := json.NewDecoder(resp1.Body).Decode(&updateScheduleParticipants); err != nil {
		return nil, err
	}

	return &updateScheduleParticipants, nil

}
