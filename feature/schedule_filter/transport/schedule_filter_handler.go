package transport

import (
	"api/feature/schedule_filter/service"
	"github.com/gofiber/fiber/v2"
	"io"
)

type ScheduleFilterHandler struct {
	service service.ScheduleFilterService
}

// ScheduleFilter godoc
// @Summary Get schedules by filter
// @Description Get schedules by filter
// @Tags Schedule
// @Accept json
// @Produce json
// @Param param query string false "Filter parameter"
// @Success 200 {array} core_dtos.TwScheduleResponse
// @Router /dbms/v1/schedule [get]
func (h *ScheduleFilterHandler) ScheduleFilter(c *fiber.Ctx) error {
	resp, err := h.service.ScheduleFilter(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response body",
		})
	}
	return c.Status(resp.StatusCode).Send(body)
}
