package transport

import (
	"api/service/document"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {array} document_dtos.TwDocumentResponse
// @Router /api/v1/document/schedule/{schedule_id} [get]
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
// @Param wspUserId formData string true "Workspace user ID who uploads the file"
// @Success 200 {string} string "File uploaded successfully"
// @Failure 400 {string} string "Bad Request - Missing or invalid parameters"
// @Failure 500 {string} string "Internal Server Error - Something went wrong during file upload"
// @Router /api/v1/document/upload [post]
func (h *DocumentHandler) uploadHandler(c *fiber.Ctx) error {

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to retrieve file"})
	}

	scheduleId := c.FormValue("scheduleId")
	if scheduleId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Must have schedule id"})
	}
	wspUserId := c.FormValue("wspUserId")
	if wspUserId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Must have workspace user id"})
	}

	bucketName := "timewise-docs"
	objectName := fmt.Sprintf("%s/%s", scheduleId, file.Filename)

	if err := h.service.UploadFileToGCS(file, bucketName, objectName, scheduleId, wspUserId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendString("File uploaded successfully to Google Cloud Storage")
}
