package transport

import (
	"api/service/schedule_log"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ScheduleLogHandler struct {
	service schedule_log.ScheduleLogService
}

func NewScheduleLogHandler() *ScheduleLogHandler {
	service := schedule_log.NewScheduleLogService()
	return &ScheduleLogHandler{
		service: *service,
	}
}

// getScheduleLogsByScheduleId godoc
// @Summary Get schedule logs by schedule ID
// @Description Get schedule logs by schedule ID
// @Tags schedule_log
// @Accept json
// @Produce json
// @Param scheduleId path string true "Schedule ID"
// @Success 200 {array} schedule_log_dtos.TwScheduleLogResponse
// @Router /api/v1/schedule_log/schedule/{scheduleId} [get]
func (h *ScheduleLogHandler) GetScheduleLogByScheduleID(c *fiber.Ctx) error {
	scheduleIDStr := c.Params("scheduleID")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid schedule ID")
	}
	scheduleLog, err := h.service.GetScheduleLogsByScheduleID(scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(scheduleLog)
}
