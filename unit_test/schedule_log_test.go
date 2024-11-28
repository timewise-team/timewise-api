package unit_test_test

import (
	"api/service/schedule_log"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc40_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_log.NewScheduleLogService()

	_, err := service.GetScheduleLogsByScheduleID("")

	assert.Equal(t, "schedule id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc40_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_log.NewScheduleLogService()

	scheduleParticipant, err := service.GetScheduleLogsByScheduleID("5")

	if err != nil {
		t.Logf("Received error: %v", err)
		t.FailNow()
	}

	assert.NoError(t, err)
	assert.Equal(t, 99, len(scheduleParticipant))
	mockDMS.AssertExpectations(t)
}

func TestFunc40_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_log.NewScheduleLogService()

	_, err := service.GetScheduleLogsByScheduleID("0")

	assert.Equal(t, "schedule id is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc40_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := schedule_log.NewScheduleLogService()

	_, err := service.GetScheduleLogsByScheduleID("-1")

	assert.Equal(t, "schedule id is required", err.Error())
	mockDMS.AssertExpectations(t)
}
