package workspace_log_service

import (
	"api/service/workspace_log"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockDMSClient struct {
	mock.Mock
}

func TestFunc47_UTCID01(t *testing.T) {
	t.Log("Func47_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_log.NewWorkspaceLogService()
	_, err := service.GetWorkspaceLogs("114")
	assert.Nil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc47_UTCID02(t *testing.T) {
	t.Log("Func47_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_log.NewWorkspaceLogService()
	_, err := service.GetWorkspaceLogs("")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc47_UTCID03(t *testing.T) {
	t.Log("Func47_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_log.NewWorkspaceLogService()
	_, err := service.GetWorkspaceLogs("30000")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc47_UTCID04(t *testing.T) {
	t.Log("Func47_UTCID04")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_log.NewWorkspaceLogService()
	_, err := service.GetWorkspaceLogs("0")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
