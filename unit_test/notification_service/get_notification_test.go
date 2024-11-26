package notification_service

import (
	noti "api/service/notfication"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockDMSClientGetNotification struct {
	mock.Mock
}

func TestFunc49_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetNotification)
	service := noti.NewNotificationService()
	userId := []string{"6", "7"}
	result, err := service.GetNotifications(userId)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockDMS.AssertExpectations(t)
}

func TestFunc49_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetNotification)
	service := noti.NewNotificationService()
	userId := []string{"6", "ahbcnbcxn"}
	_, err := service.GetNotifications(userId)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc49_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetNotification)
	service := noti.NewNotificationService()
	userId := []string{""}
	_, err := service.GetNotifications(userId)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}
