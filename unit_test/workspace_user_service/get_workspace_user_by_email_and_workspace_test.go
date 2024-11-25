package workspace_user_service

import (
	"api/service/workspace_user"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

func TestFunc16_UTCID01(t *testing.T) {
	t.Log("Func19_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	expected := &models.TwWorkspaceUser{
		ID:          202,
		WorkspaceId: 114,
		UserEmailId: 3,
		Role:        "admin",
		Status:      "joined",
		IsActive:    true,
		IsVerified:  true,
	}
	email := "giakhanh006@gmail.com"
	workspaceId := "114"
	a, err := service.GetWorkspaceUserByEmailAndWorkspaceID(email, workspaceId)

	assert.Nil(t, err)
	assert.Equal(t, expected.ID, a.ID)
	assert.Equal(t, expected.WorkspaceId, a.WorkspaceId)
	assert.Equal(t, expected.UserEmailId, a.UserEmailId)
	assert.Equal(t, expected.Role, a.Role)
	assert.Equal(t, expected.Status, a.Status)
	assert.Equal(t, expected.IsActive, a.IsActive)
	assert.Equal(t, expected.IsVerified, a.IsVerified)

	mockDMS.AssertExpectations(t)
}

func TestFunc16_UTCID02(t *testing.T) {
	t.Log("Func16_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	//expected := &models.TwWorkspaceUser{
	//	ID:          202,
	//	WorkspaceId: 114,
	//	UserEmailId: 3,
	//	Role:        "admin",
	//	Status:      "joined",
	//	IsActive:    true,
	//	IsVerified:  true,
	//}
	email := "giakhanh@gmail.com"
	workspaceId := "114"
	a, err := service.GetWorkspaceUserByEmailAndWorkspaceID(email, workspaceId)

	assert.Nil(t, err)
	assert.Equal(t, 0, a.ID)

	mockDMS.AssertExpectations(t)
}
func TestFunc16_UTCID03(t *testing.T) {
	t.Log("Func16_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	//expected := &models.TwWorkspaceUser{
	//	ID:          202,
	//	WorkspaceId: 114,
	//	UserEmailId: 3,
	//	Role:        "admin",
	//	Status:      "joined",
	//	IsActive:    true,
	//	IsVerified:  true,
	//}
	email := "giakhanh"
	workspaceId := "114"
	_, err := service.GetWorkspaceUserByEmailAndWorkspaceID(email, workspaceId)

	assert.Equal(t, "Invalid email", err.Error())

	mockDMS.AssertExpectations(t)
}

func TestFunc16_UTCID04(t *testing.T) {
	t.Log("Func16_UTCID04")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	//expected := &models.TwWorkspaceUser{
	//	ID:          202,
	//	WorkspaceId: 114,
	//	UserEmailId: 3,
	//	Role:        "admin",
	//	Status:      "joined",
	//	IsActive:    true,
	//	IsVerified:  true,
	//}
	email := "giakhanh@gmail.com"
	workspaceId := "114"
	a, err := service.GetWorkspaceUserByEmailAndWorkspaceID(email, workspaceId)

	assert.Nil(t, err)
	assert.Equal(t, 0, a.ID)

	mockDMS.AssertExpectations(t)
}
func TestFunc16_UTCID05(t *testing.T) {
	t.Log("Func19_UTCID05")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	//expected := &models.TwWorkspaceUser{
	//	ID:          202,
	//	WorkspaceId: 114,
	//	UserEmailId: 3,
	//	Role:        "admin",
	//	Status:      "joined",
	//	IsActive:    true,
	//	IsVerified:  true,
	//}
	email := "giakhanh006@gmail.com"
	workspaceId := "abcxuz"
	_, err := service.GetWorkspaceUserByEmailAndWorkspaceID(email, workspaceId)

	assert.NotNil(t, err)
	//assert.Equal(t, expected.ID, a.ID)
	//assert.Equal(t, expected.WorkspaceId, a.WorkspaceId)
	//assert.Equal(t, expected.UserEmailId, a.UserEmailId)
	//assert.Equal(t, expected.Role, a.Role)
	//assert.Equal(t, expected.Status, a.Status)
	//assert.Equal(t, expected.IsActive, a.IsActive)
	//assert.Equal(t, expected.IsVerified, a.IsVerified)

	mockDMS.AssertExpectations(t)
}

func TestFunc16_UTCID06(t *testing.T) {
	t.Log("Func19_UTCID06")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	email := "giakhanh006@gmail.com"
	workspaceId := "10000"
	_, err := service.GetWorkspaceUserByEmailAndWorkspaceID(email, workspaceId)

	assert.Equal(t, "Workspace not found", err.Error())
	//assert.Equal(t, expected.ID, a.ID)
	//assert.Equal(t, expected.WorkspaceId, a.WorkspaceId)
	//assert.Equal(t, expected.UserEmailId, a.UserEmailId)
	//assert.Equal(t, expected.Role, a.Role)
	//assert.Equal(t, expected.Status, a.Status)
	//assert.Equal(t, expected.IsActive, a.IsActive)
	//assert.Equal(t, expected.IsVerified, a.IsVerified)

	mockDMS.AssertExpectations(t)
}
func TestFunc16_UTCID07(t *testing.T) {
	t.Log("Func19_UTCID07")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace_user.NewWorkspaceUserService()
	email := "giakhanh006@gmail.com"
	workspaceId := "0"
	_, err := service.GetWorkspaceUserByEmailAndWorkspaceID(email, workspaceId)

	assert.Equal(t, "Workspace not found", err.Error())
	//assert.Equal(t, expected.ID, a.ID)
	//assert.Equal(t, expected.WorkspaceId, a.WorkspaceId)
	//assert.Equal(t, expected.UserEmailId, a.UserEmailId)
	//assert.Equal(t, expected.Role, a.Role)
	//assert.Equal(t, expected.Status, a.Status)
	//assert.Equal(t, expected.IsActive, a.IsActive)
	//assert.Equal(t, expected.IsVerified, a.IsVerified)

	mockDMS.AssertExpectations(t)
}
