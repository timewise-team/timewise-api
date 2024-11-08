package document

import (
	"api/dms"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/document_dtos"
	"github.com/timewise-team/timewise-models/models"
	"google.golang.org/api/option"
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
func (s *DocumentService) CheckIfFileExists(bucketName string, objectName string) bool {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("serviceAccount.json"))
	if err != nil {
		fmt.Printf("Failed to create client: %v", err)
		return false
	}
	defer client.Close()

	// Kiểm tra sự tồn tại của object
	_, err = client.Bucket(bucketName).Object(objectName).Attrs(ctx)
	if err != nil {
		// Nếu lỗi là "not found", thì file không tồn tại
		if err == storage.ErrObjectNotExist {
			return false
		}
		// Các lỗi khác sẽ được ghi log
		fmt.Printf("Error checking if file exists: %v", err)
		return false
	}

	// Nếu không có lỗi thì file tồn tại
	return true
}
func (s *DocumentService) UploadFileToGCS(fileHeader *multipart.FileHeader, bucketName string, objectName string, scheduleId string, wspUserId string, fileNameWithoutSchedule string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("serviceAccount.json"))
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
	wc.ContentDisposition = "attachment; filename=" + fileNameWithoutSchedule
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

	//fileNameWithoutExt := fileHeader.Filename[:len(fileHeader.Filename)-len(ext)]
	scheduleIdInt, err := strconv.Atoi(scheduleId)
	if err != nil {
		return err
	}
	wspUserIdInt, err := strconv.Atoi(wspUserId)
	if err != nil {
		return err
	}
	document := models.TwDocument{
		FileName:    fileNameWithoutSchedule,
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

func (s *DocumentService) DeleteFileFromGCS(bucketName string, objectName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("serviceAccount.json"))
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)

	// Thực hiện xóa file
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("Failed to delete file from bucket: %v", err)
	}

	return nil
}

func (s *DocumentService) DeleteDocumentFromDatabase(scheduleId string, fileName string) error {
	// Xóa bản ghi với điều kiện theo `scheduleId` và `fileName`
	resp, err := dms.CallAPI(
		"DELETE",
		"/document",
		nil,
		nil,
		map[string]string{"scheduleId": scheduleId, "fileName": fileName},
		120,
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete document in database: %v", err)
	}

	return nil
}

func (s *DocumentService) DownloadDocuments(documentId string) (*http.Response, string, error) {
	// Tìm document trong cơ sở dữ liệu
	resp, err := dms.CallAPI("GET", "/document/"+documentId, nil, nil, nil, 120)
	if err != nil {
		return nil, "", errors.New("failed to retrieve document")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", errors.New("failed to read document")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("failed to retrieve document")
	}
	var document models.TwDocument
	if err := json.Unmarshal(body, &document); err != nil {
		return nil, "", errors.New("failed to unmarshal document")
	}
	// Sử dụng FilePath hoặc DownloadUrl từ document
	downloadURL := document.DownloadUrl
	if downloadURL == "" {
		return nil, "", errors.New("no download URL available")
	}

	// Tải file từ URL đã ký
	resp, err = http.Get(downloadURL)
	if err != nil {
		return nil, "", errors.New("failed to download document")
	}

	// Kiểm tra xem file có được tải thành công không
	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("failed to retrieve document from storage")
	}

	return resp, document.FileName, nil
}
