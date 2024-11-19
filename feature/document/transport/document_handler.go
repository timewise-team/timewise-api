package transport

import (
	"api/dms"
	"api/service/document"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"io"
	"net/url"
	"path/filepath"
	"strconv"
)

type DocumentHandler struct {
	service document.DocumentService
}

func NewDocumentHandler() *DocumentHandler {
	service := document.NewDocumentService()
	return &DocumentHandler{
		service: *service,
	}
}

// GetDocumentByScheduleID godoc
// @Summary Get documents by schedule
// @Description Get documents by schedule
// @Tags document
// @Accept json
// @Produce json
// @Param scheduleId path string true "Schedule ID"
// @Security bearerToken
// @Success 200 {array} document_dtos.TwDocumentResponse
// @Router /api/v1/document/schedule/{scheduleId} [get]
func (h *DocumentHandler) GetDocumentByScheduleID(c *fiber.Ctx) error {
	scheduleIDStr := c.Params("scheduleID")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid schedule ID")
	}
	document, err := h.service.GetDocumentsByScheduleID(scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(document)
}

// @Summary Upload file
// @Description Upload file with multipart/form-data and upload to Google Cloud Storage
// @Tags document
// @Accept mpfd
// @Security bearerToken
// @Produce json
// @Param file formData file true "File to upload"
// @Param scheduleId formData string true "Schedule ID associated with the file"
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Success 200 {string} string "File uploaded successfully"
// @Failure 400 {string} string "Bad Request - Missing or invalid parameters"
// @Failure 500 {string} string "Internal Server Error - Something went wrong during file upload"
// @Router /api/v1/document/upload [post]
func (h *DocumentHandler) uploadHandler(c *fiber.Ctx) error {
	scheduleId := c.FormValue("scheduleId")
	if scheduleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Must have schedule id"})
	}
	workspaceUserLocal := c.Locals("workspace_user")
	if workspaceUserLocal == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	workspaceUser, ok := workspaceUserLocal.(*models.TwWorkspaceUser)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	wspUserId := workspaceUser.ID
	if wspUserId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not get workspace user id"})
	}
	wspUserIdStr := strconv.Itoa(wspUserId)
	// get role to check permission
	//isEditable, err := checkRole(c, scheduleId, wspUserIdStr)
	//if !isEditable {
	//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You do not have permission to upload file"})
	//}
	//if err != nil {
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	//}
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to retrieve file"})
	}
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if file.Size > maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File size exceeds the 10MB limit"})
	}
	bucketName := "timewise-docs"
	originalFileName := file.Filename
	objectName := fmt.Sprintf("%s/%s", scheduleId, originalFileName)

	// Check if file with the same name already exists and append numbering if necessary
	newFileName := originalFileName
	counter := 1
	for h.service.CheckIfFileExists(bucketName, objectName) {
		ext := filepath.Ext(originalFileName)
		baseName := originalFileName[:len(originalFileName)-len(ext)]
		newFileName = fmt.Sprintf("%s(%d)%s", baseName, counter, ext)
		objectName = fmt.Sprintf("%s/%s", scheduleId, newFileName)
		counter++
	}

	if err := h.service.UploadFileToGCS(file, bucketName, objectName, scheduleId, wspUserIdStr, newFileName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendString("File uploaded successfully to Google Cloud Storage")
}

func checkRole(c *fiber.Ctx, scheduleId string, wspUserId string) (bool, error) {
	// Get user_email_id by userId
	userId := c.Locals("userid").(string)
	resp, err := dms.CallAPI("GET", "/user_email/user/"+userId, nil, nil, nil, 120)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != fiber.StatusOK {
		return false, errors.New("error from external service: " + string(body))
	}
	var userResponse []models.TwUserEmail
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return false, errors.New("could not unmarshal response body: " + err.Error())
	}
	// parse userResponse to get list of user_email_id
	user_email_id := make([]string, len(userResponse))
	for i, user := range userResponse {
		user_email_id[i] = strconv.Itoa(user.ID)
	}

	// Get workspace IDs for the user
	resp, err = dms.CallAPI("POST", "/workspace_user/user_email_id", user_email_id, nil, nil, 120)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if resp.StatusCode != fiber.StatusOK {
		return false, errors.New("error from external service: " + string(body))
	}
	// có chứa role rồi
	var workspaceResponse []models.TwWorkspaceUser
	err = json.Unmarshal(body, &workspaceResponse)
	if err != nil {
		return false, errors.New("could not unmarshal response body: " + err.Error())
	}
	userRoles := map[int]string{} // Map workspace_id -> role
	for _, wsp := range workspaceResponse {
		userRoles[wsp.WorkspaceId] = wsp.Role
	}
	for _, wsp := range workspaceResponse {
		if wsp.Role == "admin" || wsp.Role == "owner" {
			// Admin/Owner can do anything with document of schedule
			return true, nil
		} else {
			resp, err := dms.CallAPI("GET", "/schedule/participants/"+scheduleId, nil, nil, nil, 120)
			if err != nil {
				return false, err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return false, err
			}
			if resp.StatusCode != fiber.StatusOK {
				return false, errors.New("error from external service: " + string(body))
			}
			var participants []models.TwScheduleParticipant
			err = json.Unmarshal(body, &participants)
			if err != nil {
				return false, errors.New("could not unmarshal response body: " + err.Error())
			}
			for _, participant := range participants {
				if strconv.Itoa(participant.WorkspaceUserId) == wspUserId {
					return true, nil
				}
			}
		}
	}

	return false, errors.New("no role found for the user")
}

// @Summary Delete file
// @Description Delete file from Google Cloud Storage
// @Tags document
// @Accept json
// @Security bearerToken
// @Produce json
// @Param scheduleId query string true "Schedule ID associated with the file"
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Param fileName query string true "Name of the file to delete"
// @Success 200 {string} string "File deleted successfully"
// @Failure 400 {string} string "Bad Request - Missing or invalid parameters"
// @Failure 500 {string} string "Internal Server Error - Something went wrong during file deletion"
// @Router /api/v1/document/delete [delete]
func (h *DocumentHandler) deleteHandler(c *fiber.Ctx) error {
	scheduleId := c.Query("scheduleId")
	if scheduleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Must have schedule id"})
	}
	workspaceUserLocal := c.Locals("workspace_user")
	if workspaceUserLocal == nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	workspaceUser, ok := workspaceUserLocal.(*models.TwWorkspaceUser)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	wspUserId := workspaceUser.ID
	if wspUserId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not get workspace user id"})
	}
	fileName := c.Query("fileName")
	if fileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Must have file name"})
	}

	bucketName := "timewise-docs"
	objectName := fmt.Sprintf("%s/%s", scheduleId, fileName)

	if err := h.service.DeleteFileFromGCS(bucketName, objectName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Xóa bản ghi trong database nếu cần
	if err := h.service.DeleteDocumentFromDatabase(scheduleId, fileName); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete document from database"})
	}

	return c.SendString("File deleted successfully from Google Cloud Storage")
}

// downloadDocument godoc
// @Summary Download document
// @Description Download a document from Google Cloud Storage
// @Tags document
// @Security bearerToken
// @Produce application/octet-stream
// @Param documentId path int true "Document ID"
// @Success 200 {file} file "File downloaded successfully"
// @Failure 404 {string} string "Document not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/v1/document/download/{documentId} [get]
func (h *DocumentHandler) downloadDocument(c *fiber.Ctx) error {
	// Lấy documentId từ đường dẫn
	documentID := c.Params("documentId")
	if documentID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid document ID")
	}
	// call service
	resp, fileName, err := h.service.DownloadDocuments(documentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	encodedFileName := url.PathEscape(fileName)
	// Đặt header để tải về file dưới dạng attachment
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedFileName))
	c.Set("Content-Type", "application/octet-stream")

	// Gửi file về client
	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send document"})
	}

	return c.SendString("File downloaded successfully.")
}
