package account_service

import (
	"api/service/account"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"testing"
)

type mockDMSClientUnLink struct {
	mock.Mock
}

func TestFunc10_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUnLink)
	service := account.NewAccountService()

	email := "ngkkhanh006@gmail.com"
	user := core_dtos.GetUserResponseDto{ID: 1, Email: []core_dtos.EmailDto{
		{
			Email:  email,
			Status: "",
		},
	}}
	result, err := service.UnlinkAnEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockDMS.AssertExpectations(t)
}

func TestFunc10_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUnLink)
	service := account.NewAccountService()

	email := "giakhanh"

	result, err := service.UnlinkAnEmail(email)

	assert.Error(t, err)
	assert.Equal(t, "failed to fetch email details", err.Error())
	assert.Empty(t, result)
	mockDMS.AssertExpectations(t)
}
func TestFunc10_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientUnLink)
	service := account.NewAccountService()

	email := "builanviet@gmail.com"

	result, err := service.UnlinkAnEmail(email)

	assert.Error(t, err)
	assert.Equal(t, "email is not linked to any user", err.Error())
	assert.Empty(t, result)
	mockDMS.AssertExpectations(t)
}
