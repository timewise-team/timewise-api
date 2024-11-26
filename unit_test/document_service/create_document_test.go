package document_service

import (
	"api/service/document"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"testing"
)

type mockDMSClientCreateDocs struct {
	mock.Mock
}

func TestFunc42_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientCreateDocs)
	service := document.NewDocumentService()
	fileContent := []byte("test file content")
	file := multipart.FileHeader{
		Filename: "example.txt",
		Size:     int64(len(fileContent)),
	}

	err := service.UploadFileToGCS(&file, "timewise-docs", "example.txt", "74", "94", "example.txt")

	assert.NoError(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc42_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientCreateDocs)
	service := document.NewDocumentService()
	fileContent := []byte("test file content")
	file := multipart.FileHeader{
		Filename: "example.txt",
		Size:     int64(len(fileContent)),
	}

	err := service.UploadFileToGCS(&file, "timewise-docs", "example.txt", "", "94", "example.txt")

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc42_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientCreateDocs)
	service := document.NewDocumentService()
	fileContent := []byte("test file content")
	file := multipart.FileHeader{
		Filename: "example.txt",
		Size:     int64(len(fileContent)),
	}

	err := service.UploadFileToGCS(&file, "timewise-docs", "example.txt", "74", "", "example.txt")

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}

func TestFunc42_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientCreateDocs)
	service := document.NewDocumentService()
	fileContent := []byte("test file content")
	file := multipart.FileHeader{
		Filename: "example.txt",
		Size:     int64(len(fileContent)),
	}

	err := service.UploadFileToGCS(&file, "timewise-docs", "example.txt", "74", "94", "example.txt")

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}
func TestFunc42_UTCID05(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(mockDMSClientCreateDocs)
	service := document.NewDocumentService()
	fileContent := []byte("test file content")
	file := multipart.FileHeader{
		Filename: "example.txt",
		Size:     int64(len(fileContent)),
	}

	err := service.UploadFileToGCS(&file, "timewise-docs", "example.txt", "74", "0", "example.txt")

	assert.Error(t, err)
	mockDMS.AssertExpectations(t)
}
