package notification_setting_service

import (
	"api/service/notification_setting"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockDMSClientGetNotiSetting struct {
	mock.Mock
}

func TestFunc50_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetNotiSetting)
	service := notification_setting.NewNotificationSettingService()
	userId := "6"
	result, err := service.GetNotificationSettingByUserId(userId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDMS.AssertExpectations(t)
}
