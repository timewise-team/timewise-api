package transport

import (
	"api/service/schedule"
	"api/service/schedule_participant"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

type ScheduleHandler struct {
	service schedule.ScheduleService
}

func NewScheduleHandler() *ScheduleHandler {
	service := schedule.NewScheduleService()
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
// @Security bearerToken
// @Param X-User-Email header string true "User email"
// @Param X-Workspace-Id header string true "Workspace ID"g
// @Param schedule body core_dtos.TwCreateScheduleRequest true "Schedule"
// @Success 201 {object} core_dtos.TwCreateShecduleResponse
// @Router /api/v1/schedules [post]
func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var CreateScheduleDto core_dtos.TwCreateScheduleRequest
	if err := c.BodyParser(&CreateScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	result, statusCode, err := h.service.CreateSchedule(c, CreateScheduleDto)
	if err != nil {
		if err.Error() == "permission denied" {
			return c.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(statusCode).JSON(result)
}

// GetScheduleById godoc
// @Summary Get schedule by ID
// @Description Get schedule by ID
// @Tags schedule
// @Accept json
// @Produce json
// @Security bearerToken
// @Param X-User-Email header string true "User email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Param schedule_id path int true "Schedule ID"
// @Success 200 {object} core_dtos.TwScheduleResponse
// @Router /api/v1/schedules/{schedule_id} [get]
func (h *ScheduleHandler) GetScheduleByID(c *fiber.Ctx) error {
	scheduleID := c.Params("scheduleId")
	schedule, err := h.service.GetScheduleById(scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	workspaceUser := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if workspaceUser == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	participant, err := schedule_participant.NewScheduleParticipantService().GetScheduleParticipantsByScheduleAndWorkspaceUser(scheduleID, strconv.Itoa(workspaceUser.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if workspaceUser.Role == "owner" || workspaceUser.Role == "admin" {
		schedule.ExtraData = "full-access"
	} else {
		if participant.ID == 0 {
			schedule.ExtraData = "view-only"
		} else {
			schedule.ExtraData = participant.Status
		}
	}
	return c.JSON(schedule)
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
	scheduleID := c.Params("scheduleId")
	scheduleParticipant, ok := c.Locals("scheduleParticipant").(models.TwScheduleParticipant)
	if !ok {
		return errors.New("Failed to retrieve schedule participant")
	}

	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return errors.New("Failed to retrieve workspace user")
	}
	updatedSchedule, err := h.service.UpdateSchedule(scheduleID, scheduleParticipant, workspaceUser, UpdateScheduleDto)
	if err != nil {
		if err.Error() == "permission denied" {
			return c.Status(fiber.StatusForbidden).SendString(err.Error())
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(updatedSchedule)
}

// UpdateSchedulePosition godoc
// @Summary Update an existing schedule position
// @Description Update an existing schedule position
// @Tags schedule
// @Accept json
// @Produce json
// @Param schedule_id path int true "Schedule ID"
// @Param schedule body core_dtos.TwUpdateSchedulePosition true "Schedule"
// @Success 200 {object} core_dtos.TwUpdateScheduleResponse
// @Router /api/v1/schedules/position/{schedule_id} [put]
func (h *ScheduleHandler) UpdateSchedulePosition(c *fiber.Ctx) error {
	scheduleID := c.Params("scheduleId")
	var UpdateScheduleDto core_dtos.TwUpdateSchedulePosition
	if err := c.BodyParser(&UpdateScheduleDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("invalid workspaceuser")
	}

	updateSchedule, err := h.service.UpdateSchedulePosition(scheduleID, workspaceUser, UpdateScheduleDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(updateSchedule)
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
