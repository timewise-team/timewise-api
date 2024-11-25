package board_column_service

import (
	"api/service/board_columns"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc26_UTCID01(t *testing.T) {
	t.Log("Func24=5_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	err := service.DeleteBoardColumn("0")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc26_UTCID02(t *testing.T) {
	t.Log("Func24=5_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	err := service.DeleteBoardColumn("abcxyz")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
