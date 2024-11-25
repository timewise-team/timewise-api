package workspace_user_service

import (
	"api/service/workspace_user"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

func TestFunc20_UTCID01(t *testing.T) {
	t.Log("Func20_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 51,
	}
	email := "guakhanh006@gmail.com"
	err := service.VerifyWorkspaceUserInvitation(workspaceUser, email)
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc20_UTCID02(t *testing.T) {
	t.Log("Func20_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 0,
	}
	email := "guakhanh006@gmail.com"
	err := service.VerifyWorkspaceUserInvitation(workspaceUser, email)
	assert.NotNil(t, err)
	assert.Equal(t, "workspace not found", err.Error())
	mockDMS.AssertExpectations(t)
}
func TestFunc20_UTCID03(t *testing.T) {
	t.Log("Func20_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 51,
	}
	email := ""
	err := service.VerifyWorkspaceUserInvitation(workspaceUser, email)
	assert.NotNil(t, err)
	assert.Equal(t, "email not found", err.Error())
	mockDMS.AssertExpectations(t)
}
func TestFunc20_UTCID04(t *testing.T) {
	t.Log("Func20_UTCID04")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 51,
	}
	email := "giakhanh"
	err := service.VerifyWorkspaceUserInvitation(workspaceUser, email)
	assert.NotNil(t, err)
	assert.Equal(t, "email is invalid", err.Error())
	mockDMS.AssertExpectations(t)
}
