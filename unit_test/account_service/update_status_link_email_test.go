package account_service

import (
	"api/service/account"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"testing"
)

type mockDMSClientUpdateStatusLinkEmail struct {
	mock.Mock
}

func TestFunc54_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateStatusLinkEmail)
	service := account.NewAccountService()

	userId := "2"
	email := "giakhanh@gmail.com"
	status := "pending"
	expected := core_dtos.GetUserResponseDto{
		ID:    2,
		Email: []core_dtos.EmailDto{{Email: "giakhanh@gmail.com", Status: "pending"}},
	}
	result, err := service.UpdateStatusLinkEmailRequest(email, userId, status)

	assert.NoError(t, err)
	assert.Equal(t, expected.Email, result.Email)
	mockDMS.AssertExpectations(t)
}

func TestFunc54_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateStatusLinkEmail)
	service := account.NewAccountService()

	userId := "abcbcbcb"
	email := "giakhanh@gmail.com"
	status := "pending"

	_, err := service.UpdateStatusLinkEmailRequest(email, userId, status)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc54_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateStatusLinkEmail)
	service := account.NewAccountService()

	userId := ""
	email := "giakhanh@gmail.com"
	status := "pending"

	_, err := service.UpdateStatusLinkEmailRequest(email, userId, status)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc54_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateStatusLinkEmail)
	service := account.NewAccountService()

	userId := "2"
	email := "giakhanh"
	status := "pending"

	_, err := service.UpdateStatusLinkEmailRequest(email, userId, status)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc54_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateStatusLinkEmail)
	service := account.NewAccountService()

	userId := "2"
	email := ""
	status := "pending"

	_, err := service.UpdateStatusLinkEmailRequest(email, userId, status)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc54_UTCID06(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUpdateStatusLinkEmail)
	service := account.NewAccountService()

	userId := "2"
	email := "giakhanh@gmail.com"
	status := "acvbxcvb"

	_, err := service.UpdateStatusLinkEmailRequest(email, userId, status)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}
