package workspace_user_service

import (
	"github.com/stretchr/testify/mock"
)

type MockDMSClient struct {
	mock.Mock
}

//

//func TestFunc4_UTCID02(t *testing.T) {
//	t.Log("Func4_UTCID02")
//	utils.InitConfig()
//	mockDMS := new(MockDMSClient)
//	service := workspace.NewWorkspaceService()
//
//	request := "abcxyz"
//	a := service.GetWorkspaceById(request)
//
//	assert.Nil(t, a)
//
//	mockDMS.AssertExpectations(t)
//}
//
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
