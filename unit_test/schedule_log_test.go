package unit_test_test

import (
	"api/service/schedule"
	"api/service/schedule_log"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc36_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_log.NewScheduleLogService()

	_, err := service.GetScheduleLogsByScheduleID("")

	assert.Equal(t, "schedule id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc36_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_log.NewScheduleLogService()

	scheduleParticipant, err := service.GetScheduleLogsByScheduleID("5")

	if err != nil {
		t.Logf("Received error: %v", err)
		t.FailNow()
	}

	assert.NoError(t, err)
	assert.Equal(t, 7, len(scheduleParticipant))
	mockDMS.AssertExpectations(t)
}

func TestFunc36_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetScheduleByID("0")

	assert.Equal(t, "GET /schedule/0 returned status 404: Schedule not found", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc36_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule.NewScheduleService()

	_, err := service.GetScheduleByID("999")

	assert.Equal(t, "GET /schedule/999 returned status 404: Schedule not found", err.Error())
	mockDMS.AssertExpectations(t)
	mockDMS.AssertExpectations(t)
}
