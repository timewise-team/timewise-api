package workspace_service

import (
	"api/service/workspace"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

func TestFunc4_UTCID01(t *testing.T) {
	t.Log("Func4_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	expected := models.TwWorkspace{
		ID:        29,
		Title:     "personal",
		Type:      "personal",
		IsDeleted: false,
	}
	request := "29"
	a := service.GetWorkspaceById(request)
	assert.Equal(t, expected.ID, a.ID)
	assert.Equal(t, expected.Title, a.Title)
	assert.Equal(t, expected.Type, a.Type)
	assert.Equal(t, expected.IsDeleted, a.IsDeleted)

	mockDMS.AssertExpectations(t)
}

func TestFunc4_UTCID02(t *testing.T) {
	t.Log("Func4_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	request := "abcxyz"
	a := service.GetWorkspaceById(request)

	assert.Nil(t, a)

	mockDMS.AssertExpectations(t)
}

func TestFunc4_UTCID03(t *testing.T) {
	t.Log("Func4_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	request := "36"
	a := service.GetWorkspaceById(request)

	assert.Nil(t, a)

	mockDMS.AssertExpectations(t)
}
