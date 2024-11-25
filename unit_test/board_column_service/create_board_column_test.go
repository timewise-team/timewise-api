package board_column_service

import (
	"api/service/board_columns"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/board_columns_dtos"
	"testing"
)

type MockDMSClient struct {
	mock.Mock
}

func TestFunc23_UTCID01(t *testing.T) {
	t.Log("Func23_UTCID01")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()
	request := board_columns_dtos.BoardColumnsRequest{
		WorkspaceId: 114,
		Name:        "test",
		Position:    1,
	}
	_, err := service.CreateBoardColumn(request)
	assert.NotNil(t, err)
	assert.Equal(t, "position is invalid", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc23_UTCID02(t *testing.T) {
	t.Log("Func23_UTCID02")
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := board_columns.NewBoardColumnsService()
	request := board_columns_dtos.BoardColumnsRequest{
		WorkspaceId: 2000,
		Name:        "test",
		Position:    1,
	}
	_, err := service.CreateBoardColumn(request)
	assert.NotNil(t, err)
	assert.Equal(t, "workspace not found", err.Error())
	mockDMS.AssertExpectations(t)
}
