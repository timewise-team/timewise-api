package account_service

import (
	"api/service/account"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"testing"
)

type mockDMSClientGetInfo struct {
	mock.Mock
}

func TestFunc11_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetInfo)
	service := account.NewAccountService()
	userId := "6"
	user := core_dtos.GetUserResponseDto{
		ID:    6,
		Email: []core_dtos.EmailDto{{Email: "anh.nguyenduc.work@gmail.com"}},
	}
	result, err := service.GetUserInfoByUserId(userId, "")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user, result)
	mockDMS.AssertExpectations(t)
}
func TestFunc11_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetInfo)
	service := account.NewAccountService()

	userId := "0"

	result, err := service.GetUserInfoByUserId(userId, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockDMS.AssertExpectations(t)
}
func TestFunc11_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetInfo)
	service := account.NewAccountService()

	userId := ""

	result, err := service.GetUserInfoByUserId(userId, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockDMS.AssertExpectations(t)
}
