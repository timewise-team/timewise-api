package workspace_service

import (
	"api/service/workspace"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

func TestFunc3_UTCID01(t *testing.T) {
	t.Log("Func3_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	request := "builanviet@gmail.com"
	a, err := service.GetWorkspacesByEmail(request)

	if err != nil {
		t.Logf("Received error: %v", err)
		t.FailNow()
	}
	assert.NoError(t, err)
	assert.Equal(t, []models.TwWorkspace(nil), a)
	assert.Nil(t, a)

	mockDMS.AssertExpectations(t)
}
func TestFunc3_UTCID02(t *testing.T) {
	t.Log("Func3_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	request := "abcxyz"
	a, err := service.GetWorkspacesByEmail(request)

	assert.Error(t, err)
	assert.Nil(t, a)

	mockDMS.AssertExpectations(t)
}
func TestFunc3_UTCID03(t *testing.T) {
	t.Log("Func3_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	request := "a & gmail.com"
	a, err := service.GetWorkspacesByEmail(request)

	assert.Error(t, err)
	assert.Nil(t, a)

	mockDMS.AssertExpectations(t)
}
