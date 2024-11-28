package unit_test_test

import (
	"api/service/schedule"
	"api/service/schedule_participant"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"github.com/timewise-team/timewise-models/models"
	"testing"
	"time"
)

func TestFunc37_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	scheduleService := schedule.NewScheduleService()

	_, err := scheduleService.GetScheduleByID("999")

	assert.Equal(t, "GET /schedule/999 returned status 404: Schedule not found", err.Error())

	mockDMS.AssertExpectations(t)
}

func TestFunc37_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	_, err := service.GetScheduleParticipantsByScheduleID(0)

	assert.Equal(t, "schedule not found", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc37_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	_, err := service.GetScheduleParticipantsByScheduleID(-1)

	assert.Equal(t, "schedule not found", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc37_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	scheduleParticipant, _ := service.GetScheduleParticipantsByScheduleID(100)

	assert.Equal(t, 3, len(scheduleParticipant))
	mockDMS.AssertExpectations(t)
}

func TestFunc38_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	invitedMember := schedule_participant_dtos.InviteToScheduleRequest{
		ScheduleId: 97,
		Email:      "",
	}

	_, _, err := service.InviteToSchedule(&workspaceUser, invitedMember)

	assert.Equal(t, "email is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc38_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	invitedMember := schedule_participant_dtos.InviteToScheduleRequest{
		ScheduleId: 97,
		Email:      "123123",
	}

	_, _, err := service.InviteToSchedule(&workspaceUser, invitedMember)

	assert.Equal(t, "invalid email format", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc38_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID:          2,
		WorkspaceId: 1,
	}

	invitedMember := schedule_participant_dtos.InviteToScheduleRequest{
		ScheduleId: 97,
		Email:      "quangthuan210103@gmail.com",
	}

	scheduleParticipant, _, err := service.InviteToSchedule(&workspaceUser, invitedMember)

	assert.NoError(t, err)
	assert.Equal(t, "pending", scheduleParticipant.InvitationStatus)
	mockDMS.AssertExpectations(t)
}

func TestFunc38_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID:          2,
		WorkspaceId: 1,
	}

	invitedMember := schedule_participant_dtos.InviteToScheduleRequest{
		ScheduleId: 0,
		Email:      "quangthuan210103@gmail.com",
	}

	_, _, err := service.InviteToSchedule(&workspaceUser, invitedMember)

	assert.Equal(t, "invalid schedule id", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc38_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID:          2,
		WorkspaceId: 1,
	}

	invitedMember := schedule_participant_dtos.InviteToScheduleRequest{
		ScheduleId: -1,
		Email:      "quangthuan210103@gmail.com",
	}

	_, _, err := service.InviteToSchedule(&workspaceUser, invitedMember)

	assert.Equal(t, "invalid schedule id", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc39_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	now := time.Now()

	memberAssigned := models.TwScheduleParticipant{
		ID:               131,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ScheduleId:       102,
		WorkspaceUserId:  22,
		Status:           "assign to",
		AssignAt:         &now,
		AssignBy:         2,
		ResponseTime:     &now,
		InvitationSentAt: &now,
		InvitationStatus: "joined",
	}

	member, _ := service.AssignMember(&workspaceUser, &memberAssigned)

	assert.Equal(t, "assign to", member.Status)
	mockDMS.AssertExpectations(t)
}

func TestFunc39_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	now := time.Now()

	memberAssigned := models.TwScheduleParticipant{
		ID:               131,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ScheduleId:       102,
		WorkspaceUserId:  22,
		Status:           "assign to",
		AssignAt:         &now,
		AssignBy:         2,
		ResponseTime:     &now,
		InvitationSentAt: &now,
		InvitationStatus: "pending",
	}

	_, err := service.AssignMember(&workspaceUser, &memberAssigned)

	assert.Equal(t, "member can't be assign", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc39_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	now := time.Now()

	memberAssigned := models.TwScheduleParticipant{
		ID:               131,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ScheduleId:       102,
		WorkspaceUserId:  22,
		Status:           "assign to",
		AssignAt:         &now,
		AssignBy:         2,
		ResponseTime:     &now,
		InvitationSentAt: &now,
		InvitationStatus: "declined",
	}

	_, err := service.AssignMember(&workspaceUser, &memberAssigned)

	assert.Equal(t, "member can't be assign", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc39_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	now := time.Now()

	memberAssigned := models.TwScheduleParticipant{
		ID:               131,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ScheduleId:       102,
		WorkspaceUserId:  22,
		Status:           "",
		AssignAt:         &now,
		AssignBy:         2,
		ResponseTime:     &now,
		InvitationSentAt: &now,
		InvitationStatus: "joined",
	}

	_, err := service.AssignMember(&workspaceUser, &memberAssigned)

	assert.Equal(t, "member can't be assign", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc39_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 0,
	}

	now := time.Now()

	memberAssigned := models.TwScheduleParticipant{
		ID:               131,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ScheduleId:       102,
		WorkspaceUserId:  22,
		Status:           "assign to",
		AssignAt:         &now,
		AssignBy:         2,
		ResponseTime:     &now,
		InvitationSentAt: &now,
		InvitationStatus: "joined",
	}

	_, err := service.AssignMember(&workspaceUser, &memberAssigned)

	assert.Equal(t, "invalid workspace user id", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc39_UTCID06(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_participant.NewScheduleParticipantService()

	workspaceUser := models.TwWorkspaceUser{
		ID: -1,
	}

	now := time.Now()

	memberAssigned := models.TwScheduleParticipant{
		ID:               131,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ScheduleId:       102,
		WorkspaceUserId:  22,
		Status:           "assign to",
		AssignAt:         &now,
		AssignBy:         2,
		ResponseTime:     &now,
		InvitationSentAt: &now,
		InvitationStatus: "joined",
	}

	_, err := service.AssignMember(&workspaceUser, &memberAssigned)

	assert.Equal(t, "invalid workspace user id", err.Error())
	mockDMS.AssertExpectations(t)
}
