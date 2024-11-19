package transport

import (
	"api/service/reminder"
	"api/service/schedule"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
	"time"
)

type ReminderHandler struct {
	service reminder.ReminderService
}

func NewReminderHandler() *ReminderHandler {
	service := reminder.NewReminderService()
	return &ReminderHandler{
		service: *service,
	}
}

type CreateReminderRequest struct {
	ScheduleId   int    `json:"schedule_id"`
	ReminderTime string `json:"reminder_time"`
	//type of reminder (only me, all participants)
	//Type string `json:"type"`
}

// CreateReminderAllParticipant godoc
// @Summary Create reminder for all participants
// @Description Create reminder for all participants
// @Tags reminder
// @Accept json
// @Produce json
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Param reminder body CreateReminderRequest true "Reminder"
// @Success 201 {object} models.TwReminder
// @Router /api/v1/reminder/all_participants [post]
func (h ReminderHandler) CreateReminderAllParticipant(ctx *fiber.Ctx) error {
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
	WorkspaceUser := ctx.Locals("workspace_user").(*models.TwWorkspaceUser)
	if WorkspaceUser == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	ScheduleParticipant := ctx.Locals("scheduleParticipant").(models.TwScheduleParticipant)
	if ScheduleParticipant.ID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	reminderTimeInt, err := strconv.Atoi(reminder.ReminderTime)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder time format",
		})
	}
	if reminderTimeInt < 0 || reminderTimeInt > 10080 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder time",
		})
	}
	if reminder.ScheduleId == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid schedule ID",
		})
	}
	err = h.service.CreateReminderAllParticipant(scheduleDetail, WorkspaceUser, ScheduleParticipant, reminderTimeInt)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Reminder created successfully",
	})
}

// GetRemindersAllParticipant godoc
// @Summary Get reminder for all participants
// @Description Get reminder for all participants
// @Tags reminder
// @Accept json
// @Produce json
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Param schedule_id path int true "Schedule ID"
// @Success 200 {object} models.TwReminder
// @Router /api/v1/reminder/schedule/{schedule_id}/all_participants [get]
func (h ReminderHandler) GetRemindersAllParticipant(ctx *fiber.Ctx) error {
	scheduleId := ctx.Params("scheduleId")
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
	if reminderResponse.ID == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Reminder not found",
		})
	}
	return ctx.JSON(reminderResponse)
}

type UpdateReminderRequest struct {
	ScheduleId   int    `json:"schedule_id"`
	ReminderTime string `json:"reminder_time"`
}

// UpdateReminderAllParticipant godoc
// @Summary Update reminder for all participants
// @Description Update reminder for all participants
// @Tags reminder
// @Accept json
// @Produce json
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Param reminder_id path int true "Reminder ID"
// @Param reminder body UpdateReminderRequest true "Reminder"
// @Success 200 {object} fiber.Map
// @Router /api/v1/reminder/all_participants/{reminder_id} [put]
func (h ReminderHandler) UpdateReminderAllParticipant(ctx *fiber.Ctx) error {

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
	reminderRequest, err := h.service.GetReminderByID(reminderId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get reminder",
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

	scheduleDetail, err := schedule.NewScheduleService().GetScheduleDetailByID(strconv.Itoa(reminder.ScheduleId))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get schedule detail",
		})
	}
	reminderTime := scheduleDetail.StartTime.Add(-time.Duration(reminderTimeInt) * time.Minute)

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		// Xử lý lỗi nếu không thể tải múi giờ
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error loading time zone",
		})
	}

	reminderTimeInLocal := time.Date(reminderTime.Year(), reminderTime.Month(), reminderTime.Day(), reminderTime.Hour(), reminderTime.Minute(), reminderTime.Second(), reminderTime.Nanosecond(), loc)
	fmt.Printf("reminderTime: %v\n", reminderTimeInLocal)
	fmt.Printf("time.Now(): %v\n", time.Now().In(loc))
	if reminderTimeInLocal.Before(time.Now().In(loc)) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Reminder time must be after current time",
		})
	}
	reminderRequest.ReminderTime = reminderTime
	reminderRequest.IsSent = false
	err = h.service.UpdateReminder(reminderId, reminderRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update reminder",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reminder updated successfully",
	})
}

// DeleteReminder godoc
// @Summary Delete reminder
// @Description Delete reminder
// @Tags reminder
// @Accept json
// @Produce json
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Param reminder_id path int true "Reminder ID"
// @Param schedule_id path int true "Schedule ID"
// @Success 200 {object} fiber.Map
// @Router /api/v1/reminder/{reminder_id}/schedule/{schedule_id} [delete]
func (h ReminderHandler) DeleteReminder(ctx *fiber.Ctx) error {

	reminderId := ctx.Params("reminder_id")
	if reminderId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder ID",
		})
	}
	err := h.service.DeleteReminder(reminderId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete reminder",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reminder deleted successfully",
	})
}

type CreateReminderOnlyMeRequest struct {
	ScheduleId int `json:"schedule_id"`
}

// CreateReminderOnlyMe godoc
// @Summary Create reminder for only me
// @Description Create reminder for only me
// @Tags reminder
// @Accept json
// @Produce json
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Param reminder body CreateReminderOnlyMeRequest true "Reminder"
// @Success 201 {object} fiber.Map
// @Router /api/v1/reminder/only_me [post]
func (h ReminderHandler) CreateReminderOnlyMe(ctx *fiber.Ctx) error {
	var createReminderOnlyMeRequest CreateReminderOnlyMeRequest
	if err := ctx.BodyParser(&createReminderOnlyMeRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if createReminderOnlyMeRequest.ScheduleId == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid schedule ID",
		})
	}
	scheduleIdStr := strconv.Itoa(createReminderOnlyMeRequest.ScheduleId)
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
	ScheduleParticipant := ctx.Locals("scheduleParticipant").(models.TwScheduleParticipant)
	if ScheduleParticipant.ID == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	reminderChecks, err := reminder.NewReminderService().GetRemindersByScheduleID(scheduleIdStr)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get reminders",
		})
	}
	for _, reminderCheck := range reminderChecks {
		if reminderCheck.Type == "only me" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Reminder already exists",
			})
		}
	}
	startTime := scheduleDetail.StartTime
	var reminder models.TwReminder
	if startTime != nil {
		reminder = models.TwReminder{
			ScheduleId:      createReminderOnlyMeRequest.ScheduleId,
			ReminderTime:    *startTime,
			Type:            "only me",
			IsSent:          false,
			WorkspaceUserID: WorkspaceUser.ID,
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Start time is missing",
		})
	}
	_, err = h.service.CreateReminder(reminder)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create reminder",
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Reminder created successfully",
	})

}

// GetRemindersOnlyMe godoc
// @Summary Get reminder for only me
// @Description Get reminder for only me
// @Tags reminder
// @Accept json
// @Produce json
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Param schedule_id path int true "Schedule ID"
// @Success 200 {object} models.TwReminder
// @Router /api/v1/reminder/schedule/{schedule_id}/only_me [get]
func (h ReminderHandler) GetRemindersOnlyMe(ctx *fiber.Ctx) error {
	scheduleId := ctx.Params("scheduleId")
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
		if reminder.Type == "only me" {
			reminderResponse = reminder
		}
	}
	if reminderResponse.ID == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Reminder not found",
		})
	}
	return ctx.JSON(reminderResponse)
}

type UpdateReminderOnlyMeRequest struct {
	ScheduleId int    `json:"schedule_id"`
	Time       string `json:"time"`
}

// UpdateReminderOnlyMe godoc
// @Summary Update reminder for only me
// @Description Update reminder for only me
// @Tags reminder
// @Accept json
// @Produce json
// @Param X-User-Email header string true "User Email"
// @Param X-Workspace-Id header string true "Workspace ID"
// @Security bearerToken
// @Param reminder_id path int true "Reminder ID"
// @Param reminder body UpdateReminderOnlyMeRequest true "Reminder"
// @Success 200 {object} fiber.Map
// @Router /api/v1/reminder/only_me/{reminder_id} [put]
func (h ReminderHandler) UpdateReminderOnlyMe(ctx *fiber.Ctx) error {
	reminderId := ctx.Params("reminder_id")
	if reminderId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder ID",
		})
	}

	var reminder UpdateReminderOnlyMeRequest
	if err := ctx.BodyParser(&reminder); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// reminder.Time = "2024-11-11 11:11:11"
	// Chuyển chuỗi time thành kiểu time.Time
	reminderTime, err := time.Parse("2006-01-02 15:04", reminder.Time)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder time format",
		})
	}

	if reminder.ScheduleId == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid schedule ID",
		})
	}

	if reminderTime.IsZero() {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid reminder time",
		})
	}

	reminderRequest, err := h.service.GetReminderByID(reminderId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get reminder",
		})
	}

	scheduleDetail, err := schedule.NewScheduleService().GetScheduleDetailByID(strconv.Itoa(reminder.ScheduleId))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get schedule detail",
		})
	}

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		// Xử lý lỗi nếu không thể tải múi giờ
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error loading time zone",
		})
	}

	reminderTimeInLocal := time.Date(reminderTime.Year(), reminderTime.Month(), reminderTime.Day(), reminderTime.Hour(), reminderTime.Minute(), reminderTime.Second(), reminderTime.Nanosecond(), loc)
	if reminderTimeInLocal.Before(time.Now().In(loc)) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Reminder time must be after current time",
		})
	}

	if scheduleDetail.EndTime != nil {
		if reminderTime.After(*scheduleDetail.EndTime) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Reminder time must be before schedule end time",
			})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Schedule end time is missing",
		})
	}

	reminderRequest.ReminderTime = reminderTime
	reminderRequest.IsSent = false
	err = h.service.UpdateReminder(reminderId, reminderRequest)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update reminder",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reminder updated successfully",
	})

}
