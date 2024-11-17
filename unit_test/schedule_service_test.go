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

func TestGetScheduleById_Success(t *testing.T) {
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

func TestGetScheduleById_Error(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetScheduleByID("1")

	if err != nil {
		t.Logf("Received error: %v", err)
		t.FailNow()
	}

	mockDMS.AssertExpectations(t)
}

func TestUpdateSchedule_Success(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	title := "Business"
	description := "Task for business"
	startTimeStr := "2024-11-17 16:00:00.000"
	endTimeStr := "2024-11-17 21:00:00.000"

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

func ptrString(s string) *string {
	return &s
}

func TestUpdateSchedule_Error(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	tests := []struct {
		name        string
		request     core_dtos.TwUpdateScheduleRequest
		expectErr   bool
		expectedErr string
	}{
		{
			name: "Title is empty",
			request: core_dtos.TwUpdateScheduleRequest{
				Title:       nil,
				Description: ptrString("Task for business"),
				StartTime:   ptrString("2024-11-17 16:00:00.000"),
				EndTime:     ptrString("2024-11-17 17:00:00.000"),
			},
			expectErr:   true,
			expectedErr: "title cannot be empty",
		},
		{
			name: "StartTime greater than EndTime",
			request: core_dtos.TwUpdateScheduleRequest{
				Title:       ptrString("Business"),
				Description: ptrString("Task for business"),
				StartTime:   ptrString("2024-11-17 17:00:00.000"),
				EndTime:     ptrString("2024-11-17 16:00:00.000"),
			},
			expectErr:   true,
			expectedErr: "end time cannot be earlier than start time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedulePartipant := models.TwScheduleParticipant{
				Status: "creator",
			}
			workspaceUser := models.TwWorkspaceUser{
				ID: 4,
			}

			_, err := service.UpdateSchedule("5", schedulePartipant, &workspaceUser, tt.request)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.expectedErr {
					t.Errorf("Expected error '%s', but got '%v'", tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}

	mockDMS.AssertExpectations(t)
}
