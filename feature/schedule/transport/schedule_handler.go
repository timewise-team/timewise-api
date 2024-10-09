package transport

import (
	"api/service/service"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
)

type ScheduleHandler struct {
	service service.ScheduleService
}

func NewScheduleHandler() *ScheduleHandler {
	service := service.NewScheduleService()
	return &ScheduleHandler{
		service: *service,
	}
}

// CreateSchedule godoc
// @Summary Create a new schedule
// @Description Create a new schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule body core_dtos.TwCreateScheduleRequest true "Schedule"
// @Success 201 {object} core_dtos.TwCreateShecduleResponse
// @Router /api/v1/schedules [post]
func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var CreateScheduleDto core_dtos.TwCreateScheduleRequest
	if err := c.BodyParser(&CreateScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := h.service.CreateSchedule(c, CreateScheduleDto)
	if err != nil {
		if err.Error() == "permission denied" {
			return c.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "Create schedule successfully",
	})
}

// UpdateSchedule godoc
// @Summary Update an existing schedule
// @Description Update an existing schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Param schedule body core_dtos.TwUpdateScheduleRequest true "Schedule"
// @Success 200 {object} core_dtos.TwUpdateScheduleResponse
// @Router /api/v1/schedules/{schedule_id} [put]
func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	var UpdateScheduleDto core_dtos.TwUpdateScheduleRequest
	if err := c.BodyParser(&UpdateScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := h.service.UpdateSchedule(c, UpdateScheduleDto)
	if err != nil {
		if err.Error() == "permission denied" {
			return c.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Update schedule successfully",
	})
}

// DeleteSchedule godoc
// @Summary Delete a schedule
// @Description Delete a schedule
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Success 204 "No Content"
// @Router /api/v1/schedules/{schedule_id} [delete]
func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {

	err := h.service.DeleteSchedule(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Delete schedule successfully",
	})
}
