package transport

import (
	"api/service/schedule_participant"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"strconv"
)

type ScheduleParticipantHandler struct {
	service schedule_participant.ScheduleParticipantService
}

func NewScheduleParticipantHandler() *ScheduleParticipantHandler {
	service := schedule_participant.NewScheduleParticipantService()
	return &ScheduleParticipantHandler{
		service: *service,
	}
}

// getScheduleParticipantsByScheduleId godoc
// @Summary Get schedule participants by schedule ID
// @Description Get schedule participants by schedule ID
// @Tags schedule_participant
// @Accept json
// @Produce json
// @Param scheduleId path string true "Schedule ID"
// @Success 200 {array} schedule_participant_dtos.ScheduleParticipantInfo
// @Router /api/v1/schedule_participant/schedule/{scheduleId} [get]
func (h *ScheduleParticipantHandler) GetScheduleParticipantByScheduleID(c *fiber.Ctx) error {
	scheduleIDStr := c.Params("scheduleID")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid schedule ID")
	}
	scheduleParticipant, err := h.service.GetScheduleParticipantsByScheduleID(scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(scheduleParticipant)
}

func (h *ScheduleParticipantHandler) InviteToSchedule(c *fiber.Ctx) error {
	var InviteToScheduleDto schedule_participant_dtos.InviteToScheduleRequest
	if err := c.BodyParser(&InviteToScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	scheduleParticipant, err := h.service.InviteToSchedule(c, InviteToScheduleDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(scheduleParticipant)
}
