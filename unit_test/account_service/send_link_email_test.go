package account_service

import (
	"api/service/account"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

type mockDMSClientSendLink struct {
	mock.Mock
}

func TestFunc09_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientSendLink)
	service := account.NewAccountService()

	userId := "6"
	email := "ngkkhanh006@gmail.com"
	status := "pending"
	expected := models.TwUserEmail{ID: 12, Email: email, Status: &status}

	result, err := service.SendLinkAnEmailRequest(userId, email)

	assert.NoError(t, err)
	assert.Equal(t, expected.Email, result.Email)
	mockDMS.AssertExpectations(t)
}

func TestFunc09_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientSendLink)
	service := account.NewAccountService()

	_, err := service.SendLinkAnEmailRequest("abcdez", "ngkkhanh006@gmail.com")

	assert.Equal(t, "email is already linked or pending", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc09_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientSendLink)
	service := account.NewAccountService()

	_, err := service.SendLinkAnEmailRequest("6", "giakhanh")

	assert.Equal(t, "email is not a user", err.Error())
	mockDMS.AssertExpectations(t)
}
