package workspace_user_service

import (
	"api/service/workspace_user"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

type MockDMSClient struct {
	mock.Mock
}

func TestFunc12_UTCID01(t *testing.T) {
	t.Log("Func12_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()

	workspaceId := 20
	workspaceUser := &models.TwWorkspaceUser{
		ID: 2,
	}
	a := service.AcceptInvitation(workspaceUser, workspaceId)

	assert.Equal(t, "workspace user not found", a.Error())

	mockDMS.AssertExpectations(t)
}

func TestFunc12_UTCID02(t *testing.T) {
	t.Log("Func12_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()

	workspaceId := 34
	workspaceUser := &models.TwWorkspaceUser{
		ID: 2,
	}
	a := service.AcceptInvitation(workspaceUser, workspaceId)

	assert.Equal(t, "workspace user not found", a.Error())

	mockDMS.AssertExpectations(t)
}

func TestFunc12_UTCID03(t *testing.T) {
	t.Log("Func12_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()

	workspaceId := 25
	workspaceUser := &models.TwWorkspaceUser{
		ID: 28,
	}
	a := service.AcceptInvitation(workspaceUser, workspaceId)

	assert.Equal(t, "workspace user is not pending", a.Error())

	mockDMS.AssertExpectations(t)
}
func TestFunc12_UTCID04(t *testing.T) {
	t.Log("Func12_UTCID04")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()

	workspaceId := 57
	workspaceUser := &models.TwWorkspaceUser{
		ID: 188,
	}
	a := service.AcceptInvitation(workspaceUser, workspaceId)

	assert.Nil(t, a)

	mockDMS.AssertExpectations(t)
}

//func TestFunc4_UTCID03(t *testing.T) {
//	t.Log("Func4_UTCID03")
//	utils.InitConfig()
//	mockDMS := new(MockDMSClient)
//	service := workspace.NewWorkspaceService()
//
//	request := "36"
//	a := service.GetWorkspaceById(request)
//
//	assert.Nil(t, a)
//
//	mockDMS.AssertExpectations(t)
//}
