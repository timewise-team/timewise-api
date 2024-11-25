package workspace_user_service

import (
	"api/service/workspace_user"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/workspace_user_dtos"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

func TestFunc19_UTCID01(t *testing.T) {
	t.Log("Func19_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()

	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 114,
	}
	request := workspace_user_dtos.UpdateWorkspaceUserRoleRequest{
		Email: "giakhanh006@gmail.com",
		Role:  "admin",
	}
	a := service.UpdateWorkspaceUserRole(workspaceUser, request)

	assert.Nil(t, a)

	mockDMS.AssertExpectations(t)
}

func TestFunc19_UTCID02(t *testing.T) {
	t.Log("Func19_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()

	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 114,
	}
	request := workspace_user_dtos.UpdateWorkspaceUserRoleRequest{
		Email: "builanviet",
		Role:  "admin",
	}
	a := service.UpdateWorkspaceUserRole(workspaceUser, request)

	assert.NotNil(t, a)
	assert.Equal(t, "email is invalid", a.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc19_UTCID03(t *testing.T) {
	t.Log("Func19_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()

	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 114,
	}
	request := workspace_user_dtos.UpdateWorkspaceUserRoleRequest{
		Email: "giakhanh006@gmail.com",
		Role:  "mem",
	}
	a := service.UpdateWorkspaceUserRole(workspaceUser, request)

	assert.NotNil(t, a)
	assert.Equal(t, "role is invalid", a.Error())
	mockDMS.AssertExpectations(t)
}
func TestFunc19_UTCID04(t *testing.T) {
	t.Log("Func19_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()

	workspaceUser := &models.TwWorkspaceUser{
		WorkspaceId: 114,
	}
	request := workspace_user_dtos.UpdateWorkspaceUserRoleRequest{
		Email: "t@gmail.com",
		Role:  "member",
	}
	a := service.UpdateWorkspaceUserRole(workspaceUser, request)

	assert.NotNil(t, a)
	assert.Equal(t, "workspace user not found", a.Error())
	mockDMS.AssertExpectations(t)
}
