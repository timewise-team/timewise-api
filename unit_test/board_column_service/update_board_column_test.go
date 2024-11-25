package board_column_service

import (
	"api/service/board_columns"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc25_UTCID01(t *testing.T) {
	t.Log("Func25_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	_, err := service.UpdateBoardColumn("144", "test")
	assert.Nil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc25_UTCID02(t *testing.T) {
	t.Log("Func25_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	_, err := service.UpdateBoardColumn("144", "")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc25_UTCID03(t *testing.T) {
	t.Log("Func25_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	_, err := service.UpdateBoardColumn("144", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc25_UTCID04(t *testing.T) {
	t.Log("Func25_UTCID03")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()

	_, err := service.UpdateBoardColumn("abcdef", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	assert.NotNil(t, err)
	mockDMS.AssertExpectations(t)
}
