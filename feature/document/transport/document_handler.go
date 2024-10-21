package transport

import (
	"api/service/document"
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

// getDocumentsBySchedule godoc
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
	document, err := h.service.GetDocumentsBySchedule(scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(document)
}
