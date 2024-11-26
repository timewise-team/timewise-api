package document_service

import (
	"api/service/document"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockDMSClientDownloadDocs struct {
	mock.Mock
}

func TestFunc43_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDownloadDocs)
	service := document.NewDocumentService()
	documentId := "20"
	_, fileName, err := service.DownloadDocuments(documentId)

	assert.NoError(t, err)
	assert.NotNil(t, fileName)
	mockDMS.AssertExpectations(t)
}

func TestFunc43_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDownloadDocs)
	service := document.NewDocumentService()
	documentId := "0"
	_, _, err := service.DownloadDocuments(documentId)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc43_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientDownloadDocs)
	service := document.NewDocumentService()
	documentId := "abcxnbcnx"
	_, _, err := service.DownloadDocuments(documentId)

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}
