package workspace_service

import (
	"api/service/workspace"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/models"
	"testing"
	"time"
)

func TestFunc7_UTCID01(t *testing.T) {
	t.Log("Func7_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	request := "3"
	a, err := service.GetWorkspacesByUserId(request)
	assert.NoError(t, err)
	assert.Equal(t, []models.TwWorkspace{
		{
			ID:          29,
			CreatedAt:   time.Date(2024, time.October, 9, 0, 18, 38, 859000000, time.UTC),
			UpdatedAt:   time.Date(2024, time.October, 9, 0, 18, 38, 859000000, time.UTC),
			DeletedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			Title:       "personal",
			ExtraData:   "",
			Description: "",
			Key:         "",
			Type:        "personal",
			IsDeleted:   false}, models.TwWorkspace{
			ID:          30,
			CreatedAt:   time.Date(2024, time.October, 9, 8, 23, 21, 689000000, time.UTC),
			UpdatedAt:   time.Date(2024, time.October, 9, 8, 23, 21, 689000000, time.UTC),
			DeletedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			Title:       "string",
			ExtraData:   "",
			Description: "string",
			Key:         "",
			Type:        "workspace",
			IsDeleted:   false}, models.TwWorkspace{
			ID:          32,
			CreatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			DeletedAt:   time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			Title:       "hehe",
			ExtraData:   "",
			Description: "",
			Key:         "",
			Type:        "workspace",
			IsDeleted:   false}}, a)

	mockDMS.AssertExpectations(t)
}

func TestFunc7_UTCID02(t *testing.T) {
	t.Log("Func7_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	//expected := models.TwWorkspace{
	//	ID:        29,
	//	Title:     "personal",
	//	Type:      "personal",
	//	IsDeleted: false,
	//}
	request := "abcxyz"
	_, err := service.GetWorkspacesByUserId(request)
	assert.Error(t, err)
	assert.Equal(t, "Invalid user id", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc7_UTCID03(t *testing.T) {
	t.Log("Func7_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := workspace.NewWorkspaceService()

	//expected := models.TwWorkspace{
	//	ID:        29,
	//	Title:     "personal",
	//	Type:      "personal",
	//	IsDeleted: false,
	//}
	request := "1"
	_, err := service.GetWorkspacesByUserId(request)
	assert.Error(t, err)
	assert.Equal(t, "User not found", err.Error())
	mockDMS.AssertExpectations(t)
}
