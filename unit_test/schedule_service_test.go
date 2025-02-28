package unit_test_test

import (
	"api/service/schedule"
	"api/unit_test/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"testing"
	"time"
)

func TestFunc31_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	workspaceId := 1
	boardColumnId := 4
	workspaceUserId := 2

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := core_dtos.TwCreateScheduleRequest{
		Title:           &title,
		Description:     &description,
		BoardColumnID:   &boardColumnId,
		WorkspaceID:     &workspaceId,
		WorkspaceUserID: &workspaceUserId,
	}

	response, status, err := service.CreateSchedule(&workspaceUser, request)

	assert.NoError(t, err)
	assert.Equal(t, title, response.Title)
	assert.Equal(t, 201, status)
	assert.Equal(t, description, response.Description)
	assert.Equal(t, boardColumnId, response.BoardColumnID)
	assert.Equal(t, workspaceId, response.WorkspaceID)

	mockDMS.AssertExpectations(t)
}

func TestFunc31_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	workspaceId := 1
	workspaceUserId := 2

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := core_dtos.TwCreateScheduleRequest{
		Title:           &title,
		Description:     &description,
		BoardColumnID:   nil,
		WorkspaceID:     &workspaceId,
		WorkspaceUserID: &workspaceUserId,
	}

	_, _, err := service.CreateSchedule(&workspaceUser, request)
	assert.Equal(t, "Invalid board column id", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc31_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	workspaceId := 1
	boardColumnId := 0
	workspaceUserId := 2

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := core_dtos.TwCreateScheduleRequest{
		Title:           &title,
		Description:     &description,
		BoardColumnID:   &boardColumnId,
		WorkspaceID:     &workspaceId,
		WorkspaceUserID: &workspaceUserId,
	}

	_, _, err := service.CreateSchedule(&workspaceUser, request)
	assert.Equal(t, "Invalid board column id", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc31_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	boardColumnId := 4
	workspaceUserId := 2

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := core_dtos.TwCreateScheduleRequest{
		Title:           &title,
		Description:     &description,
		BoardColumnID:   &boardColumnId,
		WorkspaceID:     nil,
		WorkspaceUserID: &workspaceUserId,
	}

	_, _, err := service.CreateSchedule(&workspaceUser, request)
	assert.Equal(t, "Invalid workspace id", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc31_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	workspaceId := 0
	boardColumnId := 4
	workspaceUserId := 2

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := core_dtos.TwCreateScheduleRequest{
		Title:           &title,
		Description:     &description,
		BoardColumnID:   &boardColumnId,
		WorkspaceID:     &workspaceId,
		WorkspaceUserID: &workspaceUserId,
	}

	_, _, err := service.CreateSchedule(&workspaceUser, request)
	assert.Equal(t, "Invalid workspace id", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc31_UTCID06(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := ""
	description := "Task for business"
	workspaceId := 1
	boardColumnId := 4
	workspaceUserId := 2

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := core_dtos.TwCreateScheduleRequest{
		Title:           &title,
		Description:     &description,
		BoardColumnID:   &boardColumnId,
		WorkspaceID:     &workspaceId,
		WorkspaceUserID: &workspaceUserId,
	}

	_, _, err := service.CreateSchedule(&workspaceUser, request)
	assert.Equal(t, "Invalid title", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc32_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetScheduleByID("")

	assert.Equal(t, "schedule id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc32_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	mockSchedule := models.TwSchedule{
		ID:          5,
		Title:       "Business",
		Description: "Task for business",
	}

	scheduleDto, err := service.GetScheduleByID("5")

	if err != nil {
		t.Logf("Received error: %v", err)
		t.FailNow()
	}

	assert.NoError(t, err)
	assert.Equal(t, mockSchedule.Title, scheduleDto.Title)
	assert.Equal(t, 5, scheduleDto.ID)
	assert.Equal(t, "Task for business", scheduleDto.Description)

	mockDMS.AssertExpectations(t)
}

func TestFunc32_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetScheduleByID("0")

	assert.Equal(t, "GET /schedule/0 returned status 404: Schedule not found", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc32_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetScheduleByID("999")

	assert.Equal(t, "GET /schedule/999 returned status 404: Schedule not found", err.Error())
	mockDMS.AssertExpectations(t)
	mockDMS.AssertExpectations(t)
}

func TestFunc33_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := ""
	description := "Task for business"
	startTimeStr := "2024-11-17 16:00:00.000"
	endTimeStr := "2024-11-17 21:00:00.000"

	schedulePartipant := models.TwScheduleParticipant{
		Status: "creator",
	}

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	request := core_dtos.TwUpdateScheduleRequest{
		Title:       &title,
		Description: &description,
		StartTime:   &startTimeStr,
		EndTime:     &endTimeStr,
	}

	_, err := service.UpdateSchedule("5", schedulePartipant, &workspaceUser, request)

	assert.Equal(t, "Bad Request: title cannot be empty", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc33_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	startTimeStr := "2024-11-17 16:00:00.000"
	endTimeStr := "2024-11-17 21:00:00.000"

	schedulePartipant := models.TwScheduleParticipant{
		Status: "creator",
	}

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	request := core_dtos.TwUpdateScheduleRequest{
		Title:       &title,
		Description: &description,
		StartTime:   &startTimeStr,
		EndTime:     &endTimeStr,
	}

	_, err := service.UpdateSchedule("5", schedulePartipant, &workspaceUser, request)

	assert.Equal(t, "Bad Request: start time cannot be in the past", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc33_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	startTimeStr := "2024-12-30 16:00:00.000"
	endTimeStr := "2024-12-30 16:00:00.000"

	schedulePartipant := models.TwScheduleParticipant{
		Status: "creator",
	}

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	request := core_dtos.TwUpdateScheduleRequest{
		Title:       &title,
		Description: &description,
		StartTime:   &startTimeStr,
		EndTime:     &endTimeStr,
	}

	_, err := service.UpdateSchedule("5", schedulePartipant, &workspaceUser, request)

	assert.Equal(t, "Bad Request: Invalid Endtime", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc33_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	startTimeStr := "2024-12-30 16:00:00.000"
	endTimeStr := "2024-12-30 14:00:00.000"

	schedulePartipant := models.TwScheduleParticipant{
		Status: "creator",
	}

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	request := core_dtos.TwUpdateScheduleRequest{
		Title:       &title,
		Description: &description,
		StartTime:   &startTimeStr,
		EndTime:     &endTimeStr,
	}

	_, err := service.UpdateSchedule("5", schedulePartipant, &workspaceUser, request)

	assert.Equal(t, "Bad Request: Invalid Endtime", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc33_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	startTimeStr := "2024-12-20 16:00:00.000"
	endTimeStr := "2024-12-20 21:00:00.000"

	startTime, err := time.Parse("2006-01-02 15:04:05.000", startTimeStr)
	if err != nil {
		fmt.Println("Error parsing start time:", err)
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05.000", endTimeStr)
	if err != nil {
		fmt.Println("Error parsing end time:", err)
		return
	}

	schedulePartipant := models.TwScheduleParticipant{
		Status: "creator",
	}

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	request := core_dtos.TwUpdateScheduleRequest{
		Title:       &title,
		Description: &description,
		StartTime:   &startTimeStr,
		EndTime:     &endTimeStr,
	}

	updatedSchedule := models.TwSchedule{
		ID:          5,
		Title:       title,
		Description: description,
		StartTime:   &startTime,
		EndTime:     &endTime,
	}

	scheduleDto, err := service.UpdateSchedule("5", schedulePartipant, &workspaceUser, request)

	assert.NoError(t, err)
	assert.Equal(t, updatedSchedule.Title, scheduleDto.Title)
	assert.Equal(t, 5, scheduleDto.ID)
	assert.Equal(t, updatedSchedule.Description, scheduleDto.Description)
	assert.Equal(t, updatedSchedule.StartTime, scheduleDto.StartTime)
	assert.Equal(t, updatedSchedule.EndTime, scheduleDto.EndTime)

	mockDMS.AssertExpectations(t)
}

func TestFunc33_UTCID06(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	startTimeStr := "2024-12-30 16:00:00.000"
	endTimeStr := "2024-12-30 18:00:00.000"

	schedulePartipant := models.TwScheduleParticipant{
		Status: "creator",
	}

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	request := core_dtos.TwUpdateScheduleRequest{
		Title:       &title,
		Description: &description,
		StartTime:   &startTimeStr,
		EndTime:     &endTimeStr,
	}

	_, err := service.UpdateSchedule("", schedulePartipant, &workspaceUser, request)

	assert.Equal(t, "schedule id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc33_UTCID07(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	startTimeStr := "2024-12-30 16:00:00.000"
	endTimeStr := "2024-12-30 18:00:00.000"

	schedulePartipant := models.TwScheduleParticipant{
		Status: "creator",
	}

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	request := core_dtos.TwUpdateScheduleRequest{
		Title:       &title,
		Description: &description,
		StartTime:   &startTimeStr,
		EndTime:     &endTimeStr,
	}

	_, err := service.UpdateSchedule("0", schedulePartipant, &workspaceUser, request)

	assert.Equal(t, "schedule id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc34_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	err := service.DeleteSchedule("", &workspaceUser)

	assert.Equal(t, "schedule not found", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc34_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	err := service.DeleteSchedule("93", &workspaceUser)

	assert.Equal(t, nil, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc34_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	err := service.DeleteSchedule("0", &workspaceUser)

	assert.Equal(t, "schedule not found", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc34_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	workspaceUser := models.TwWorkspaceUser{
		ID: 4,
	}

	err := service.DeleteSchedule("-1", &workspaceUser)

	assert.Equal(t, "schedule not found", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc36_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	schedules, _ := service.GetSchedulesByBoardColumn("1", 1)

	assert.Equal(t, 4, len(schedules))
	mockDMS.AssertExpectations(t)
}

func TestFunc36_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetSchedulesByBoardColumn("1", 0)

	assert.Equal(t, "board column id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc36_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetSchedulesByBoardColumn("1", -1)

	assert.Equal(t, "board column id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc36_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetSchedulesByBoardColumn("", -1)

	assert.Equal(t, "workspace id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc36_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetSchedulesByBoardColumn("0", -1)

	assert.Equal(t, "workspace id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc35_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	position := 4
	boardColumnId := 4

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	updatePositionDto := core_dtos.TwUpdateSchedulePosition{
		Position:      &position,
		BoardColumnID: &boardColumnId,
	}

	schedules, _ := service.UpdateSchedulePosition("102", &workspaceUser, updatePositionDto)

	assert.Equal(t, 4, schedules.BoardColumnID)
	assert.Equal(t, 4, schedules.Position)
	mockDMS.AssertExpectations(t)
}

func TestFunc35_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	position := 0
	boardColumnId := 4

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	updatePositionDto := core_dtos.TwUpdateSchedulePosition{
		Position:      &position,
		BoardColumnID: &boardColumnId,
	}

	_, err := service.UpdateSchedulePosition("102", &workspaceUser, updatePositionDto)

	assert.Equal(t, "invalid position", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc35_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	position := -1
	boardColumnId := 4

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	updatePositionDto := core_dtos.TwUpdateSchedulePosition{
		Position:      &position,
		BoardColumnID: &boardColumnId,
	}

	_, err := service.UpdateSchedulePosition("102", &workspaceUser, updatePositionDto)

	assert.Equal(t, "invalid position", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc35_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	position := 5
	boardColumnId := 0

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	updatePositionDto := core_dtos.TwUpdateSchedulePosition{
		Position:      &position,
		BoardColumnID: &boardColumnId,
	}

	_, err := service.UpdateSchedulePosition("102", &workspaceUser, updatePositionDto)

	assert.Equal(t, "invalid board column id", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc35_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	position := 5
	boardColumnId := -1

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	updatePositionDto := core_dtos.TwUpdateSchedulePosition{
		Position:      &position,
		BoardColumnID: &boardColumnId,
	}

	_, err := service.UpdateSchedulePosition("102", &workspaceUser, updatePositionDto)

	assert.Equal(t, "invalid board column id", err.Error())
	mockDMS.AssertExpectations(t)
}
