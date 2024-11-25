package workspace_user_service

import (
	"api/service/workspace_user"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

func TestFunc14_UTCID01(t *testing.T) {
	t.Log("Func14_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		ID:          202,
		WorkspaceId: 114,
	}
	workspaceId := "202"
	err := service.DeleteWorkspaceUser(workspaceUser, workspaceId)

	assert.Nil(t, err)

	mockDMS.AssertExpectations(t)
}
func TestFunc14_UTCID02(t *testing.T) {
	t.Log("Func14_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		ID:          202,
		WorkspaceId: 0,
	}
	workspaceId := "202"
	err := service.DeleteWorkspaceUser(workspaceUser, workspaceId)

	assert.NotNil(t, err)

	mockDMS.AssertExpectations(t)
}
func TestFunc14_UTCID03(t *testing.T) {
	t.Log("Func14_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		ID:          202,
		WorkspaceId: 0,
	}
	workspaceId := ""
	err := service.DeleteWorkspaceUser(workspaceUser, workspaceId)

	assert.NotNil(t, err)

	mockDMS.AssertExpectations(t)
}

func TestFunc14_UTCID04(t *testing.T) {
	t.Log("Func14_UTCID04")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		ID:          202,
		WorkspaceId: 0,
	}
	workspaceId := "abcxzy"
	err := service.DeleteWorkspaceUser(workspaceUser, workspaceId)

	assert.NotNil(t, err)

	mockDMS.AssertExpectations(t)
}
func TestFunc14_UTCID06(t *testing.T) {
	t.Log("Func14_UTCID06")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		ID:          202,
		WorkspaceId: 0,
	}
	workspaceId := "50000"
	err := service.DeleteWorkspaceUser(workspaceUser, workspaceId)

	assert.NotNil(t, err)

	mockDMS.AssertExpectations(t)
}
