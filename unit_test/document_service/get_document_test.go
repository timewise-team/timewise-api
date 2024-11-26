package document_service

import (
	"api/service/document"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockDMSClientGetDocs struct {
	mock.Mock
}

func TestFunc45_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetDocs)
	service := document.NewDocumentService()

	scheduleId := 74
	_, err := service.GetDocumentsByScheduleID(scheduleId)

	assert.NoError(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc45_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetDocs)
	service := document.NewDocumentService()

	scheduleId := -1
	_, err := service.GetDocumentsByScheduleID(scheduleId)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc45_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientGetDocs)
	service := document.NewDocumentService()

	scheduleId := 100000
	_, err := service.GetDocumentsByScheduleID(scheduleId)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}
