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

func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var CreateScheduleDto core_dtos.TwCreateScheduleRequest
	if err := c.BodyParser(&CreateScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := h.service.CreateSchedule(c, CreateScheduleDto)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Create schedule successfully",
	})
}

func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	var UpdateScheduleDto core_dtos.TwUpdateScheduleRequest
	if err := c.BodyParser(&UpdateScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := h.service.UpdateSchedule(c, UpdateScheduleDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Update schedule successfully",
	})
}

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
