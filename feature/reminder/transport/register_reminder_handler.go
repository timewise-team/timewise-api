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

	router.Post("/", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to"}), reminderHandler.Handler.CreateReminder)
	router.Get("/schedule/:schedule_id", reminderHandler.Handler.GetReminders)
	router.Put("/:reminder_id", reminderHandler.Handler.UpdateReminder)
	router.Delete("/:reminder_id", reminderHandler.Handler.DeleteReminder)
}
