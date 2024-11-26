package account_service

import (
	"api/service/account"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"testing"
)

type mockDMSClientUpdateInfo struct {
	mock.Mock
}

func TestFunc52_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateInfo)
	service := account.NewAccountService()
	userId := "2"
	firstName := "Việt"
	lastName := "Bùi"
	profilePicture := ""
	notificationSettings := ""
	calendarSettings := ""
	request := core_dtos.UpdateProfileRequestDto{
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       profilePicture,
		NotificationSettings: notificationSettings,
		CalendarSettings:     calendarSettings,
	}
	user := core_dtos.GetUserResponseDto{
		ID:                   2,
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       profilePicture,
		NotificationSettings: notificationSettings,
		CalendarSettings:     calendarSettings,
	}
	result, err := service.UpdateUserInfo(userId, request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user, result)
	mockDMS.AssertExpectations(t)
}

func TestFunc52_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateInfo)
	service := account.NewAccountService()
	userId := "abcxyz"
	firstName := "Việt"
	lastName := "Bùi"
	profilePicture := ""
	notificationSettings := ""
	calendarSettings := ""
	request := core_dtos.UpdateProfileRequestDto{
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       profilePicture,
		NotificationSettings: notificationSettings,
		CalendarSettings:     calendarSettings,
	}
	result, err := service.UpdateUserInfo(userId, request)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockDMS.AssertExpectations(t)
}

func TestFunc52_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateInfo)
	service := account.NewAccountService()
	userId := "0"
	firstName := "Việt"
	lastName := "Bùi"
	profilePicture := ""
	notificationSettings := ""
	calendarSettings := ""
	request := core_dtos.UpdateProfileRequestDto{
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       profilePicture,
		NotificationSettings: notificationSettings,
		CalendarSettings:     calendarSettings,
	}
	result, err := service.UpdateUserInfo(userId, request)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockDMS.AssertExpectations(t)
}

func TestFunc52_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateInfo)
	service := account.NewAccountService()
	userId := ""
	firstName := "Việt"
	lastName := "Bùi"
	profilePicture := ""
	notificationSettings := ""
	calendarSettings := ""
	request := core_dtos.UpdateProfileRequestDto{
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       profilePicture,
		NotificationSettings: notificationSettings,
		CalendarSettings:     calendarSettings,
	}
	result, err := service.UpdateUserInfo(userId, request)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockDMS.AssertExpectations(t)
}

func TestFunc52_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateInfo)
	service := account.NewAccountService()
	userId := "2"
	firstName := ""
	lastName := ""
	profilePicture := ""
	notificationSettings := ""
	calendarSettings := ""
	request := core_dtos.UpdateProfileRequestDto{
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       profilePicture,
		NotificationSettings: notificationSettings,
		CalendarSettings:     calendarSettings,
	}
	user := core_dtos.GetUserResponseDto{
		ID:                   2,
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       profilePicture,
		NotificationSettings: notificationSettings,
		CalendarSettings:     calendarSettings,
	}
	result, err := service.UpdateUserInfo(userId, request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user, result)
	mockDMS.AssertExpectations(t)
}

func TestFunc52_UTCID06(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateInfo)
	service := account.NewAccountService()
	userId := "2"
	firstName := ""
	lastName := "Bùi"
	profilePicture := ""
	notificationSettings := ""
	calendarSettings := ""
	request := core_dtos.UpdateProfileRequestDto{
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       profilePicture,
		NotificationSettings: notificationSettings,
		CalendarSettings:     calendarSettings,
	}
	user := core_dtos.GetUserResponseDto{
		ID:                   2,
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       profilePicture,
		NotificationSettings: notificationSettings,
		CalendarSettings:     calendarSettings,
	}
	result, err := service.UpdateUserInfo(userId, request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user, result)
	mockDMS.AssertExpectations(t)
}
