package schedule_participant

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"github.com/timewise-team/timewise-models/models"
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

func (h *ScheduleParticipantService) InviteToSchedule(c *fiber.Ctx, InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest) (*models.TwScheduleParticipant, error) {
	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return nil, errors.New("Failed to retrieve schedule participant")
	}

	resp, err := dms.CallAPI(
		"GET",
		"/workspace_user/email/"+InviteToScheduleDto.Email+"/workspace/"+strconv.Itoa(workspaceUser.WorkspaceId),
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var WorkspaceUser models.TwWorkspaceUser
	if err := json.NewDecoder(resp.Body).Decode(&WorkspaceUser); err != nil {
		return nil, err
	}

	scheduleParticipant, ok := c.Locals("scheduleParticipant").(models.TwScheduleParticipant)
	if !ok {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve schedule participant",
		})
	}

	newScheduleParticipant := models.TwScheduleParticipant{
		ScheduleId:       scheduleParticipant.ScheduleId,
		WorkspaceUserId:  WorkspaceUser.ID,
		AssignBy:         workspaceUser.ID,
		InvitationSentAt: time.Now(),
		InvitationStatus: "pending",
	}

	resp1, err := dms.CallAPI(
		"POST",
		"/schedule_participant",
		newScheduleParticipant,
		nil,
		nil,
		120,
	)

	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	var scheduleParticipants models.TwScheduleParticipant
	if err := json.NewDecoder(resp1.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	return &scheduleParticipants, nil

}

func (h *ScheduleParticipantService) AcceptInvite(c *fiber.Ctx, InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest) (*models.TwScheduleParticipant, error) {

	scheduleId := c.Params("schedule_id")
	workspaceUserId := c.Params("workspace_user_id")

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

	updateScheduleParticipant := models.TwScheduleParticipant{
		CreatedAt:        scheduleParticipants.CreatedAt,
		UpdatedAt:        time.Now(),
		ScheduleId:       scheduleParticipants.ScheduleId,
		WorkspaceUserId:  scheduleParticipants.WorkspaceUserId,
		Status:           "participant",
		AssignAt:         scheduleParticipants.AssignAt,
		AssignBy:         scheduleParticipants.AssignBy,
		ResponseTime:     time.Now(),
		InvitationSentAt: scheduleParticipants.InvitationSentAt,
		InvitationStatus: "joined",
	}

	resp1, err := dms.CallAPI(
		"PUT",
		"/schedule_participant/"+strconv.Itoa(scheduleParticipants.ID),
		updateScheduleParticipant,
		nil,
		nil,
		120,
	)

	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	var updateScheduleParticipants models.TwScheduleParticipant
	if err := json.NewDecoder(resp1.Body).Decode(&updateScheduleParticipants); err != nil {
		return nil, err
	}

	return &updateScheduleParticipants, nil

}

func (h *ScheduleParticipantService) DeclineInvite(c *fiber.Ctx, InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest) (*models.TwScheduleParticipant, error) {

	scheduleId := c.Params("schedule_id")
	workspaceUserId := c.Params("workspace_user_id")

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

	updateScheduleParticipant := models.TwScheduleParticipant{
		CreatedAt:        scheduleParticipants.CreatedAt,
		UpdatedAt:        time.Now(),
		ScheduleId:       scheduleParticipants.ScheduleId,
		WorkspaceUserId:  scheduleParticipants.WorkspaceUserId,
		Status:           scheduleParticipants.Status,
		AssignAt:         scheduleParticipants.AssignAt,
		AssignBy:         scheduleParticipants.AssignBy,
		ResponseTime:     time.Now(),
		InvitationSentAt: scheduleParticipants.InvitationSentAt,
		InvitationStatus: "declined",
	}

	resp1, err := dms.CallAPI(
		"PUT",
		"/schedule_participant/"+strconv.Itoa(scheduleParticipants.ID),
		updateScheduleParticipant,
		nil,
		nil,
		120,
	)

	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	var updateScheduleParticipants models.TwScheduleParticipant
	if err := json.NewDecoder(resp1.Body).Decode(&updateScheduleParticipants); err != nil {
		return nil, err
	}

	return &updateScheduleParticipants, nil

}
