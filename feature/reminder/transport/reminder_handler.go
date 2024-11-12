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
	ScheduleId   int `json:"schedule_id"`
	ReminderTime string `json:"reminder_time"`
	//type of reminder (only me, all participants)
	//Type string `json:"type"`
}

func (h ReminderHandler) CreateReminderAllParticipant(ctx *fiber.Ctx) error {
	var reminder CreateReminderRequest
	if err := ctx.BodyParser(&reminder); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	reminderTimeInt, err := strconv.Atoi(reminder.ReminderTime)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder time format",
		})
	}
	if reminderTimeInt < 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder time",
		})
	}
	if reminder.ScheduleId == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid schedule ID",
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
	WorkspaceUser := ctx.Locals("workspace_user").(*models.TwWorkspaceUser)
	if WorkspaceUser == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	ScheduleParticipant := ctx.Locals("scheduleParticipant").(*models.TwScheduleParticipant)
	if ScheduleParticipant == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	startTime := scheduleDetail.StartTime
	reminderTime := startTime.Add(-time.Duration(reminderTimeInt) * time.Minute)

	var reminderRequests []models.TwReminder
	if WorkspaceUser.Role == "admin" || WorkspaceUser.Role == "owner" || (WorkspaceUser.Role == "member" && ScheduleParticipant.Status == "creator") || (WorkspaceUser.Role == "member" && ScheduleParticipant.Status == "assign to") {
		//if reminder.Type == "only me" {
		var reminderRequest = models.TwReminder{
			ScheduleId:      reminder.ScheduleId,
			ReminderTime:    reminderTime,
			Type:            "all participants",
			Method:      reminder.ReminderTime,
			WorkspaceUserID: WorkspaceUser.ID,
			IsSent:          false,
		}
		reminderRequests = append(reminderRequests, reminderRequest)
	} else {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden",
		})
	}
	for _, reminderRequest := range reminderRequests {
		_, err := h.service.CreateReminder(reminderRequest)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create reminder",
			})
		}

	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Reminder created successfully",
	})
}

//	type CreateReminderRequest struct {
//		ScheduleId   int       `json:"schedule_id"`
//		Time int `json:"time"`
//	}
func (h ReminderHandler) GetRemindersAllParticipant(ctx *fiber.Ctx) error {
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
	var reminderResponse models.TwReminder
	for _, reminder := range reminders {
		if reminder.Type == "all participants" {
			reminderResponse = reminder
		}
	}
	return ctx.JSON(reminderResponse)
}

type UpdateReminderRequest struct {
	ReminderTime time.Time `json:"reminder_time"`
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
