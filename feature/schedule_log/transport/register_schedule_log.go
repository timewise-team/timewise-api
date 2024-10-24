package transport

import "github.com/gofiber/fiber/v2"

type ScheduleLogHandlerRegister struct {
	Router  fiber.Router
	Handler *ScheduleLogHandler
}

func RegisterScheduleLogHandler(router fiber.Router) {
	handler := NewScheduleLogHandler()
	scheduleLogHandler := &ScheduleLogHandlerRegister{
		Handler: handler,
	}

	// Register all endpoints here
	router.Get("schedule/:scheduleId",
		scheduleLogHandler.Handler.GetScheduleLogByScheduleID)
}
