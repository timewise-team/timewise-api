package transport

import "github.com/gofiber/fiber/v2"

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
	router.Post("/create_schedule", scheduleHandler.Handler.CreateSchedule)
	router.Put("/:scheduleId", scheduleHandler.Handler.UpdateSchedule)
	router.Delete("/:scheduleId", scheduleHandler.Handler.DeleteSchedule)
}
