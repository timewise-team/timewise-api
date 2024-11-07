package transport

import (
	"api/service/document"
	"api/service/schedule"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
	"time"
)

type ReminderHandler struct {
	service document.DocumentService
}

func NewReminderHandler() *ReminderHandler {
	service := document.NewDocumentService()
	return &ReminderHandler{
		service: *service,
	}
}

type CreateReminderRequest struct {
	ScheduleId   int       `json:"schedule_id"`
	ReminderTime time.Time `json:"reminder_time"`
	Method       string    `json:"method"`
	Type         string    `json:"type"`
}

func (h ReminderHandler) CreateReminder(ctx *fiber.Ctx) error {
	var reminder CreateReminderRequest
	if err := ctx.BodyParser(&reminder); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	scheduleIdStr := strconv.Itoa(reminder.ScheduleId)
	if scheduleIdStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid schedule ID",
		})
	}
	scheduleDetail, err := schedule.NewScheduleService().GetScheduleDetailByID(scheduleIdStr)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get schedule detail",
		})
	}
	if scheduleDetail.ID == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Schedule not found",
		})
	}
	if scheduleDetail.EndTime.Before(reminder.ReminderTime) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Reminder time must be before schedule end time",
		})
	}
	WorkspaceUser := ctx.Locals("workspace_user").(*models.TwWorkspaceUser)
	var reminderRequest = models.TwReminder{
		ScheduleId:      reminder.ScheduleId,
		ReminderTime:    reminder.ReminderTime,
		Method:          reminder.Method,
		Type:            reminder.Type,
		WorkspaceUserID: WorkspaceUser.ID,
	}
	result, err := h.service.CreateReminder(reminderRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create reminder",
		})
	}
	if result.Type == "only me" {

	} else if result.Type == "all participants" {

	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Reminder created successfully",
	})
}

func (h ReminderHandler) GetReminders(ctx *fiber.Ctx) error {
	scheduleId := ctx.Params("schedule_id")
	if scheduleId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid schedule ID",
		})
	}
	reminders, err := h.service.GetRemindersByScheduleID(scheduleId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get reminders",
		})
	}
	return ctx.JSON(reminders)
}

type UpdateReminderRequest struct {
	ReminderTime time.Time `json:"reminder_time"`
	Method       string    `json:"method"`
	Type         string    `json:"type"`
}

func (h ReminderHandler) UpdateReminder(ctx *fiber.Ctx) error {

	reminderId := ctx.Params("reminder_id")
	if reminderId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder ID",
		})
	}
	var reminder UpdateReminderRequest
	if err := ctx.BodyParser(&reminder); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	var reminderRequest UpdateReminderRequest
	reminderRequest := h.service.GetReminderByID(reminderId)
	if result := h.service.UpdateReminder(reminderId, reminder); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update reminder",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reminder updated successfully",
	})
}

func (h ReminderHandler) DeleteReminder(ctx *fiber.Ctx) error {

	reminderId := ctx.Params("reminder_id")
	if reminderId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder ID",
		})
	}
	if result := h.service.DeleteReminder(reminderId); result.Error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete reminder",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reminder deleted successfully",
	})
}
