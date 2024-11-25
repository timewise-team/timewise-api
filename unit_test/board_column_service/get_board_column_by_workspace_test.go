package board_column_service

import (
	"api/service/board_columns"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc24_UTCID01(t *testing.T) {
	t.Log("Func24_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	_, err := service.GetBoardColumnsByWorkspace("114")
	assert.Nil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc24_UTCID02(t *testing.T) {
	t.Log("Func24_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	_, err := service.GetBoardColumnsByWorkspace("114")
	assert.Nil(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc24_UTCID03(t *testing.T) {
	t.Log("Func24_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	_, err := service.GetBoardColumnsByWorkspace("")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc24_UTCID04(t *testing.T) {
	t.Log("Func24_UTCID04")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	_, err := service.GetBoardColumnsByWorkspace("abcxyz")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc24_UTCID05(t *testing.T) {
	t.Log("Func24_UTCID05")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	_, err := service.GetBoardColumnsByWorkspace("3000")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
