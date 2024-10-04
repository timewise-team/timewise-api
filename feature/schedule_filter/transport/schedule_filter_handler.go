package transport

import (
	"api/service/service"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos"
	"io"
)

type ScheduleFilterHandler struct {
	service service.ScheduleFilterService
}

func NewScheduleFilterHandler() *ScheduleFilterHandler {
	service := service.NewScheduleFilterService()
	return &ScheduleFilterHandler{
		service: *service,
	}
}

// ScheduleFilter godoc
// @Summary Get schedules by filter
// @Description Get schedules by filter
// @Tags Schedule
// @Accept json
// @Produce json
// @Param param query string false "Filter parameter"
// @Success 200 {array} dtos.TwScheduleResponse
// @Router /api/v1/schedule [get]
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
	var scheduleResponse dtos.TwScheduleResponse
	err = json.Unmarshal(body, &scheduleResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not marshal response body"})
	}
	return c.Status(resp.StatusCode).Send(body)
}
