package document_service

import (
	"api/service/document"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockDMSClientDeleteDocs struct {
	mock.Mock
}

func TestFunc44_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDeleteDocs)
	service := document.NewDocumentService()

	scheduleId := "74"
	fileName := "example.txt"
	err := service.DeleteFileFromGCS("timewise-docs", fileName)
	err2 := service.DeleteDocumentFromDatabase(scheduleId, fileName)

	assert.NoError(t, err)
	assert.NoError(t, err2)
	mockDMS.AssertExpectations(t)
}

func TestFunc44_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDeleteDocs)
	service := document.NewDocumentService()

	scheduleId := ""
	fileName := "example.txt"
	//err := service.DeleteFileFromGCS("timewise-docs", fileName)
	err2 := service.DeleteDocumentFromDatabase(scheduleId, fileName)

	//assert.Error(t, err)
	assert.Error(t, err2)
	mockDMS.AssertExpectations(t)
}

func TestFunc44_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDeleteDocs)
	service := document.NewDocumentService()

	scheduleId := "0"
	fileName := "example.txt"
	//err := service.DeleteFileFromGCS("timewise-docs", fileName)
	err2 := service.DeleteDocumentFromDatabase(scheduleId, fileName)

	//assert.NoError(t, err)
	assert.Error(t, err2)
	mockDMS.AssertExpectations(t)
}

func TestFunc44_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDeleteDocs)
	service := document.NewDocumentService()

	scheduleId := "74"
	fileName := ""
	err := service.DeleteFileFromGCS("timewise-docs", fileName)
	err2 := service.DeleteDocumentFromDatabase(scheduleId, fileName)

	assert.Error(t, err)
	assert.Error(t, err2)
	mockDMS.AssertExpectations(t)
}

func TestFunc44_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDeleteDocs)
	service := document.NewDocumentService()

	scheduleId := "74"
	fileName := "not_existed.txt"
	err := service.DeleteFileFromGCS("timewise-docs", fileName)
	err2 := service.DeleteDocumentFromDatabase(scheduleId, fileName)

	assert.Error(t, err)
	assert.Error(t, err2)
	mockDMS.AssertExpectations(t)
}
