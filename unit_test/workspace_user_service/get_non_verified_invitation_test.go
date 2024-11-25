package workspace_user_service

import (
	"api/service/workspace_user"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc17_UTCID01(t *testing.T) {
	t.Log("Func17_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceId := "114"
	_, err := service.GetWorkspaceUserInvitationNotVerifiedList(workspaceId)
	assert.Nil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc17_UTCID02(t *testing.T) {
	t.Log("Func17_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceId := ""
	_, err := service.GetWorkspaceUserInvitationNotVerifiedList(workspaceId)
	assert.NotNil(t, err)
	assert.Equal(t, "workspace id not found", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc17_UTCID03(t *testing.T) {
	t.Log("Func17_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceId := "abcxyz"
	_, err := service.GetWorkspaceUserInvitationNotVerifiedList(workspaceId)
	assert.NotNil(t, err)
	assert.Equal(t, "workspace id is invalid", err.Error())
	mockDMS.AssertExpectations(t)
}
func TestFunc17_UTCID04(t *testing.T) {
	t.Log("Func17_UTCID04")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	workspaceId := "3000"
	_, err := service.GetWorkspaceUserInvitationNotVerifiedList(workspaceId)
	assert.NotNil(t, err)
	assert.Equal(t, "workspace not found", err.Error())
	mockDMS.AssertExpectations(t)
}
