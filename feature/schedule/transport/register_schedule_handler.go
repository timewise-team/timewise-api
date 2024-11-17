package transport

import (
	"api/middleware"
	"github.com/gofiber/fiber/v2"
)

type ScheduleHandlerRegister struct {
	Router  fiber.Router
	Handler *ScheduleHandler
}

func RegisterScheduleHandler(router fiber.Router) {
	handler := NewScheduleHandler()
	scheduleHandler := &ScheduleHandlerRegister{
		Handler: handler,
	}

	// Register all endpoints here
	router.Get("/:scheduleId",
		scheduleHandler.Handler.GetScheduleByID)
	router.Post("/",
		middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		scheduleHandler.Handler.CreateSchedule)
	router.Put("/:scheduleId",
		middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to"}),
		scheduleHandler.Handler.UpdateSchedule)
	router.Delete("/:scheduleId",
		middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator"}),
		scheduleHandler.Handler.DeleteSchedule)
	router.Put("/position/:scheduleId",
		middleware.CheckWorkspaceRole([]string{"owner", "admin"}),
		scheduleHandler.Handler.UpdateSchedulePosition)
}
