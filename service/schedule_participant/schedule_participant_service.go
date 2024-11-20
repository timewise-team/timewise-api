package schedule_participant

import (
	"api/config"
	"api/dms"
	"api/notification"
	"api/service/auth"
	"api/service/user_email"
	"api/service/workspace_user"
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

func (h *ScheduleParticipantService) GetScheduleParticipantsByScheduleAndWorkspaceUser(scheduleId, workspaceId string) (*models.TwScheduleParticipant, error) {
	if scheduleId == "" {
		return nil, nil
	}
	if workspaceId == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/workspace_user/"+workspaceId+"/schedule/"+scheduleId,
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

	return &scheduleParticipants, nil

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
	InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest, check int,
) (*schedule_participant_dtos.ScheduleParticipantResponse, string, error) {

	workspaceUserInvite, err := h.getWorkspaceUserFromContext(c)
	if err != nil {
		return nil, "", err
	}

	workspaceUserInvited, err := h.GetWorkspaceUserByEmail(
		InviteToScheduleDto.Email, workspaceUserInvite.WorkspaceId,
	)
	if err != nil {
		return nil, "", err
	}

	//scheduleParticipant, err := h.getScheduleParticipantFromContext(c)
	//if err != nil {
	//	return nil, err
	//}

	schedule, err := h.getScheduleById(InviteToScheduleDto.ScheduleId)
	if err != nil {
		return nil, "", err
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, "", fiber.NewError(500, "Failed to load config")
	}

	acceptLink, declineLink, _ := auth.GenerateInviteScheduleLinks(
		cfg, InviteToScheduleDto.ScheduleId, workspaceUserInvited.ID,
	)

	scheduleParticipantResponse, err := h.handleInvitation(
		cfg, schedule, InviteToScheduleDto.ScheduleId, workspaceUserInvite, workspaceUserInvited,
		acceptLink, declineLink, InviteToScheduleDto.Email,
	)

	if err != nil {
		return nil, "", err
	}

	// create json of link
	link := map[string]string{
		"accept":  acceptLink,
		"decline": declineLink,
	}
	linkJson, _ := json.Marshal(link)

	// send notification
	notificationDto := models.TwNotifications{
		Title:       fmt.Sprintf("Invitation to join schedule %s", schedule.Title),
		Description: fmt.Sprintf("You have been invited to join schedule %s", schedule.Title),
		Link:        string(linkJson),
		UserEmailId: workspaceUserInvited.ID,
		Type:        "schedule_invitation",
	}
	err = notification.PushNotifications(notificationDto)
	if err != nil {
		return nil, "", err
	}

	return scheduleParticipantResponse, acceptLink, nil
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
func (h *ScheduleParticipantService) GetWorkspaceUserByEmail(email string, workspaceId int) (*models.TwWorkspaceUser, error) {
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
	cfg *config.Config, schedule *models.TwSchedule, scheduleId int,
	workspaceUserInvite *models.TwWorkspaceUser, workspaceUserInvited *models.TwWorkspaceUser,
	acceptLink, declineLink, email string,
) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	resp, err := dms.CallAPI(
		"GET", fmt.Sprintf("/schedule_participant/workspace_user/%d/schedule/%d",
			workspaceUserInvited.ID, scheduleId),
		nil, nil, nil, 120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	now := time.Now()
	if resp.StatusCode == http.StatusNotFound {
		return h.createAndSendInvitation(
			cfg, schedule, scheduleId, workspaceUserInvite, workspaceUserInvited,
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
	cfg *config.Config, schedule *models.TwSchedule, scheduleId int,
	workspaceUserInvite *models.TwWorkspaceUser, workspaceUserInvited *models.TwWorkspaceUser,
	acceptLink, declineLink, email string, now time.Time,
) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	newParticipant := models.TwScheduleParticipant{
		CreatedAt:        now,
		UpdatedAt:        now,
		ScheduleId:       scheduleId,
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

func (h *ScheduleParticipantService) InviteOutsideWorkspace(
	c *fiber.Ctx,
	workspaceUserInvite models.TwWorkspaceUser,
	scheduleParticipantInvite models.TwScheduleParticipant,
	InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest,
) (*models.TwWorkspaceUser, *schedule_participant_dtos.ScheduleParticipantResponse, string, error) {

	// Lấy thông tin email của người dùng
	userEmail, errs := user_email.NewUserEmailService().GetUserEmail(InviteToScheduleDto.Email)
	if userEmail == nil {
		return nil, nil, "", c.Status(500).JSON(fiber.Map{
			"message": "This email is not registered",
		})
	}
	if errs != nil {
		return nil, nil, "", c.Status(500).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	var workspaceUserResponse *models.TwWorkspaceUser
	var scheduleParticipant *schedule_participant_dtos.ScheduleParticipantResponse
	var err error
	var AcceptLink string

	// Kiểm tra vai trò của người dùng và thêm vào workspace
	if workspaceUserInvite.Role == "admin" || workspaceUserInvite.Role == "owner" {
		var temp models.TwWorkspaceUser
		temp, err = workspace_user.NewWorkspaceUserService().
			AddWorkspaceUserViaScheduleInvitation(userEmail, workspaceUserInvite.WorkspaceId, true)
		workspaceUserResponse = &temp

		// Mời người dùng tham gia lịch trình
		scheduleParticipantt, acceptLink, err := h.InviteToSchedule(c, InviteToScheduleDto, 0)
		scheduleParticipant = scheduleParticipantt
		AcceptLink = acceptLink
		if err != nil {
			return nil, nil, "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	} else if workspaceUserInvite.Role == "member" && scheduleParticipantInvite.Status == "creator" {
		var temp models.TwWorkspaceUser
		temp, err = workspace_user.NewWorkspaceUserService().
			AddWorkspaceUserViaScheduleInvitation(userEmail, workspaceUserInvite.WorkspaceId, false)
		workspaceUserResponse = &temp

		now := time.Now()
		newParticipant := models.TwScheduleParticipant{
			CreatedAt:        now,
			UpdatedAt:        now,
			ScheduleId:       scheduleParticipantInvite.ScheduleId,
			WorkspaceUserId:  workspaceUserResponse.ID,
			AssignBy:         workspaceUserInvite.ID,
			InvitationSentAt: &now,
			InvitationStatus: "pending",
		}

		resp, err := dms.CallAPI("POST", "/schedule_participant", newParticipant, nil, nil, 120)
		if err != nil {
			return nil, nil, "", err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipant); err != nil {
			return nil, nil, "", err
		}
	}

	if err != nil {
		return nil, nil, "", c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return workspaceUserResponse, scheduleParticipant, AcceptLink, nil
}

func (h *ScheduleParticipantService) AssignMember(
	c *fiber.Ctx,
	memberAssigned *models.TwScheduleParticipant,
) (*schedule_participant_dtos.ScheduleParticipantResponse, error) {

	workspaceUserInvite, err := h.getWorkspaceUserFromContext(c)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	updateScheduleParticipant := models.TwScheduleParticipant{
		CreatedAt:        memberAssigned.CreatedAt,
		UpdatedAt:        time.Now(),
		ScheduleId:       memberAssigned.ScheduleId,
		WorkspaceUserId:  memberAssigned.WorkspaceUserId,
		Status:           "assign to",
		AssignAt:         &now,
		AssignBy:         workspaceUserInvite.ID,
		ResponseTime:     memberAssigned.ResponseTime,
		InvitationSentAt: memberAssigned.InvitationSentAt,
		InvitationStatus: memberAssigned.InvitationStatus,
	}

	resp, err := dms.CallAPI(
		"PUT", fmt.Sprintf("/schedule_participant/%d", memberAssigned.ID),
		updateScheduleParticipant, nil, nil, 120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateScheduleParticipants schedule_participant_dtos.ScheduleParticipantResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateScheduleParticipants); err != nil {
		return nil, err
	}

	return &updateScheduleParticipants, nil
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

	resp2, err2 := dms.CallAPI(
		"GET",
		"/workspace_user/"+workspaceUserId,
		nil,
		nil,
		nil,
		120,
	)
	if err2 != nil {
		return nil, err2
	}
	defer resp2.Body.Close()

	var workspaceUser models.TwWorkspaceUser
	if errParsing := json.NewDecoder(resp2.Body).Decode(&workspaceUser); errParsing != nil {
		return nil, errParsing
	}

	now := time.Now()
	scheduleParticipants.Status = "participant"
	scheduleParticipants.ResponseTime = &now
	scheduleParticipants.InvitationStatus = "joined"

	workspaceUser.Status = "joined"
	workspaceUser.IsActive = true
	workspaceUser.IsVerified = true

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

	resp3, err3 := dms.CallAPI(
		"PUT",
		"/workspace_user/"+strconv.Itoa(workspaceUser.ID),
		workspaceUser,
		nil,
		nil,
		120,
	)

	if err3 != nil {
		return nil, err3
	}
	defer resp3.Body.Close()

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

	resp2, err2 := dms.CallAPI(
		"GET",
		"/workspace_user/"+workspaceUserId,
		nil,
		nil,
		nil,
		120,
	)
	if err2 != nil {
		return nil, err2
	}
	defer resp2.Body.Close()

	var workspaceUser models.TwWorkspaceUser
	if errParsing := json.NewDecoder(resp2.Body).Decode(&workspaceUser); errParsing != nil {
		return nil, errParsing
	}

	now := time.Now()
	scheduleParticipants.ResponseTime = &now
	scheduleParticipants.InvitationStatus = "declined"

	workspaceUser.Status = "declined"
	workspaceUser.IsActive = false
	workspaceUser.IsVerified = true

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

	resp3, err3 := dms.CallAPI(
		"PUT",
		"/workspace_user/"+strconv.Itoa(workspaceUser.ID),
		workspaceUser,
		nil,
		nil,
		120,
	)

	if err3 != nil {
		return nil, err3
	}
	defer resp3.Body.Close()

	return &updateScheduleParticipants, nil

}
