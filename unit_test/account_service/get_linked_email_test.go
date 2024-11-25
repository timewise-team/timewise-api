package account_service

import (
	"api/service/account"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"testing"
)

type MockDMSClient struct {
	mock.Mock
}

func TestFunc08_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	userId := "6"
	status := ""
	// expected data:
	expectedEmails := []core_dtos.EmailDto{
		{
			Email:  "anh.nguyenduc.work@gmail.com",
			Status: "",
		},
		{
			Email:  "anhndhe170145@fpt.edu.vn",
			Status: "linked",
		},
		{
			Email:  "anh.nguyenduc4@vti.com.vn",
			Status: "pending",
		},
	}
	emails, err := service.GetLinkedUserEmails(userId, status)

	assert.NoError(t, err)
	assert.NotNil(t, emails)

	for _, expectedEmail := range expectedEmails {
		found := false
		for _, email := range emails {
			if email.Email == expectedEmail.Email {
				assert.Equal(t, expectedEmail.Status, email.Status)
				found = true
				break
			}
		}
		assert.True(t, found, "Email not found: %s", expectedEmail.Email)
	}

	mockDMS.AssertExpectations(t)
}

func TestFunc08_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	userId := "6"
	status := "pending"
	// expected data:
	expectedEmails := []core_dtos.EmailDto{
		{
			Email:  "anh.nguyenduc4@vti.com.vn",
			Status: "pending",
		},
	}
	emails, err := service.GetLinkedUserEmails(userId, status)

	assert.NoError(t, err)
	assert.NotNil(t, emails)

	for _, expectedEmail := range expectedEmails {
		found := false
		for _, email := range emails {
			if email.Email == expectedEmail.Email {
				assert.Equal(t, expectedEmail.Status, email.Status)
				found = true
				break
			}
		}
		assert.True(t, found, "Email not found: %s", expectedEmail.Email)
	}

	mockDMS.AssertExpectations(t)
}

func TestFunc08_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	userId := "6"
	status := "linked"
	// expected data:
	expectedEmails := []core_dtos.EmailDto{
		{
			Email:  "anhndhe170145@fpt.edu.vn",
			Status: "linked",
		},
	}
	emails, err := service.GetLinkedUserEmails(userId, status)

	assert.NoError(t, err)
	assert.NotNil(t, emails)

	for _, expectedEmail := range expectedEmails {
		found := false
		for _, email := range emails {
			if email.Email == expectedEmail.Email {
				assert.Equal(t, expectedEmail.Status, email.Status)
				found = true
				break
			}
		}
		assert.True(t, found, "Email not found: %s", expectedEmail.Email)
	}

	mockDMS.AssertExpectations(t)
}

func TestFunc08_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	userId := "6"
	status := "abcxyz"

	emails, err := service.GetLinkedUserEmails(userId, status)

	assert.NoError(t, err)
	assert.Empty(t, emails)
	mockDMS.AssertExpectations(t)
}

func TestFunc08_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	userID := "abcxyz"

	emails, err := service.GetLinkedUserEmails(userID, "")

	assert.Error(t, err)
	assert.Empty(t, emails)
	mockDMS.AssertExpectations(t)
}

func TestFunc08_UTCID06(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := account.NewAccountService()

	userID := ""
	status := ""
	emails, err := service.GetLinkedUserEmails(userID, status)

	assert.Error(t, err)
	assert.Empty(t, emails)
	mockDMS.AssertExpectations(t)
}
