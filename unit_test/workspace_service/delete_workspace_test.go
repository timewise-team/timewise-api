package workspace_service

import (
	"api/service/workspace"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc5_UTCID01(t *testing.T) {
	t.Log("Func5_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	request := "36"
	err := service.DeleteWorkspace(request)

	assert.NoError(t, err)

	mockDMS.AssertExpectations(t)
}

//	func TestFunc5_UTCID02(t *testing.T) {
//		t.Log("Func5_UTCID02")
//		utils.InitConfig()
//		mockDMS := new(MockDMSClient)
//		service := workspace.NewWorkspaceService()
//
//		request := "36"
//		err := service.DeleteWorkspace(request)
//
//		assert.NoError(t, err)
//
//		mockDMS.AssertExpectations(t)
//	}
func TestFunc5_UTCID03(t *testing.T) {
	t.Log("Func5_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	request := "abcxyz"
	err := service.DeleteWorkspace(request)

	assert.Error(t, err)
	assert.Equal(t, "Invalid workspace id", err.Error())

	mockDMS.AssertExpectations(t)
}
func TestFunc5_UTCID04(t *testing.T) {
	t.Log("Func5_UTCID04")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	request := "36"
	err := service.DeleteWorkspace(request)

	assert.Error(t, err)

	mockDMS.AssertExpectations(t)
}
