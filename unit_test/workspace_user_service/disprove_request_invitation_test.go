package workspace_user_service

import (
	"api/service/workspace_user"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

func TestFunc25_UTCID01(t *testing.T) {
	t.Log("Func25_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 51,
	}
	email := "guakhanh006@gmail.com"
	err := service.DisproveWorkspaceUserInvitation(workspaceUser, email)
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc25_UTCID02(t *testing.T) {
	t.Log("Func25_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 51,
	}
	email := ""
	err := service.DisproveWorkspaceUserInvitation(workspaceUser, email)
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc25_UTCID03(t *testing.T) {
	t.Log("Func25_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 51,
	}
	email := "builanviet"
	err := service.DisproveWorkspaceUserInvitation(workspaceUser, email)
	assert.NotNil(t, err)
	assert.Equal(t, "email is invalid", err.Error())
	mockDMS.AssertExpectations(t)
}
func TestFunc25_UTCID04(t *testing.T) {
	t.Log("Func25_UTCID04")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 0,
	}
	email := "guakhanh006@gmail.com"
	err := service.DisproveWorkspaceUserInvitation(workspaceUser, email)
	assert.NotNil(t, err)
	assert.Equal(t, "workspace not found", err.Error())
	mockDMS.AssertExpectations(t)
}
