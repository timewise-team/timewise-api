package transport

import (
	"github.com/gofiber/fiber/v2"
)

type ScheduleFilterHandlerRegister struct {
	Router  fiber.Router
	Handler *ScheduleFilterHandler
}

func RegisterScheduleFilterHandler(router fiber.Router) {
	handler := NewScheduleFilterHandler()
	scheduleFilterHandler := &ScheduleFilterHandlerRegister{
		Handler: handler,
	}

	// Register all endpoints here
	router.Get("/schedule", scheduleFilterHandler.Handler.ScheduleFilter)
}
