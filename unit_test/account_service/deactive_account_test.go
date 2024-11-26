package account_service

import (
	"api/service/account"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockDMSClientDeactivate struct {
	mock.Mock
}

func TestFunc53_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDeactivate)
	service := account.NewAccountService()
	userId := "2"

	err := service.DeactivateAccount(userId)

	assert.NoError(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc53_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDeactivate)
	service := account.NewAccountService()
	userId := "abcbxbcx"

	err := service.DeactivateAccount(userId)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc53_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDeactivate)
	service := account.NewAccountService()
	userId := ""

	err := service.DeactivateAccount(userId)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}
