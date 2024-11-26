package notification_setting_service

import (
	"api/service/notification_setting"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

type mockDMSClientUpdateNotiSetting struct {
	mock.Mock
}

func TestFunc49_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateNotiSetting)
	service := notification_setting.NewNotificationSettingService()
	userId := "6"
	notificationSetting := models.TwNotificationSettings{
		NotificationOnTag:            true,
		NotificationOnDueDate:        true,
		NotificationOnComment:        true,
		NotificationOnScheduleChange: true,
	}
	result, err := service.UpdateNotificationSetting(userId, notificationSetting)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, notificationSetting, result)
	mockDMS.AssertExpectations(t)
}

func TestFunc49_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateNotiSetting)
	service := notification_setting.NewNotificationSettingService()
	userId := "abcxyz"
	notificationSetting := models.TwNotificationSettings{
		NotificationOnTag:            true,
		NotificationOnDueDate:        true,
		NotificationOnComment:        true,
		NotificationOnScheduleChange: true,
	}
	_, err := service.UpdateNotificationSetting(userId, notificationSetting)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc49_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateNotiSetting)
	service := notification_setting.NewNotificationSettingService()
	userId := "999990"
	notificationSetting := models.TwNotificationSettings{
		NotificationOnTag:            true,
		NotificationOnDueDate:        true,
		NotificationOnComment:        true,
		NotificationOnScheduleChange: true,
	}
	_, err := service.UpdateNotificationSetting(userId, notificationSetting)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}
