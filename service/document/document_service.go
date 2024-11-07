package document

import (
	"api/dms"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/document_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type DocumentService struct {
}

func NewDocumentService() *DocumentService {
	return &DocumentService{}
}
func (h *DocumentService) GetDocumentsBySchedule(scheduleId int) ([]models.TwDocument, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/document/schedule/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var documents []models.TwDocument
	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return nil, err
	}

	return documents, nil
}

func (h *DocumentService) GetDocumentsByScheduleID(scheduleId int) ([]document_dtos.TwDocumentResponse, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/document/schedule_id/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var documents []document_dtos.TwDocumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return nil, err
	}

	return documents, nil
}

func (s *DocumentService) UploadFileToGCS(fileHeader *multipart.FileHeader, bucketName string, objectName string, scheduleId string, wspUserId string) error {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}
	defer client.Close()

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	wc := obj.NewWriter(ctx)
	// Đặt Content-Disposition: attachment để yêu cầu tải về
	wc.ContentDisposition = "attachment; filename=" + fileHeader.Filename
	defer wc.Close()

	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("Failed to write to bucket: %v", err)
	}

	// Đặt quyền truy cập công khai cho file đã upload
	//if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
	//	return fmt.Errorf("Failed to set ACL: %v", err)
	//}
	// Lấy Signed URL thay vì sử dụng URL công khai
	signedURL, err := bucket.SignedURL(objectName, &storage.SignedURLOptions{
		Expires: time.Now().Add(24 * time.Hour),
		Method:  "GET",
	})
	if err != nil {
		return fmt.Errorf("Failed to generate signed URL: %v", err)
	}
	ext := filepath.Ext(fileHeader.Filename)
	extWithoutDot := ext[1:]

	fileNameWithoutExt := fileHeader.Filename[:len(fileHeader.Filename)-len(ext)]
	scheduleIdInt, err := strconv.Atoi(scheduleId)
	if err != nil {
		return err
	}
	wspUserIdInt, err := strconv.Atoi(wspUserId)
	if err != nil {
		return err
	}
	document := models.TwDocument{
		FileName:    fileNameWithoutExt,
		FilePath:    fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName),
		FileSize:    int(fileHeader.Size),
		FileType:    extWithoutDot,
		DownloadUrl: signedURL,
		ScheduleId:  scheduleIdInt,
		UploadedBy:  wspUserIdInt,
		UploadedAt:  time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	err = s.CreateDocumentInDatabase(document)
	if err != nil {
		return fmt.Errorf("failed to create document in database: %v", err)
	}

	return nil

	// push notification
}

func (s *DocumentService) CreateDocumentInDatabase(document models.TwDocument) error {
	resp, err := dms.CallAPI(
		"POST",
		"/document/upload",
		document,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create document in database: %v", err)
	}
	return nil
}
