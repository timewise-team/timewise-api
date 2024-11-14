package transport

import (
	"api/middleware"
	"github.com/gofiber/fiber/v2"
)

type ReminderHandlerRegister struct {
	Router  fiber.Router
	Handler *ReminderHandler
}

func RegisterReminderHandler(router fiber.Router) {
	handler := NewReminderHandler()
	reminderHandler := &ReminderHandlerRegister{
		Handler: handler,
	}

	router.Post("/all_participants", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to"}), reminderHandler.Handler.CreateReminderAllParticipant)
	router.Get("/schedule/:scheduleId/all_participants", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member", "guest"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to", "participant"}), reminderHandler.Handler.GetRemindersAllParticipant)
	router.Put("/all_participants/:reminder_id", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to"}), reminderHandler.Handler.UpdateReminderAllParticipant)
	router.Delete("/:reminder_id/schedule/:scheduleId", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to"}), reminderHandler.Handler.DeleteReminder)
	router.Post("/only_me", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member", "guest"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to", "participant"}), reminderHandler.Handler.CreateReminderOnlyMe)
	router.Get("/schedule/:scheduleId/only_me", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member", "guest"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to", "participant"}), reminderHandler.Handler.GetRemindersOnlyMe)
	router.Put("/only_me/:reminder_id", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member", "guest"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to", "participant"}), reminderHandler.Handler.UpdateReminderOnlyMe)
}
