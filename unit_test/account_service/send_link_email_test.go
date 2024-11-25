package account_service

import (
	"api/service/account"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

func TestFunc09_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
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
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	_, err := service.SendLinkAnEmailRequest("abcdez", "ngkkhanh006@gmail.com")

	assert.Equal(t, "strconv.Atoi: parsing 'abcdez': invalid syntax", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc09_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	_, err := service.SendLinkAnEmailRequest("6", "giakhanh")

	assert.Equal(t, "Email not found", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc09_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	userId := "6"
	email := "giakhanh"
	expectedError := "email is already linked or pending"

	_, err := service.SendLinkAnEmailRequest(userId, email)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	mockDMS.AssertExpectations(t)
}
