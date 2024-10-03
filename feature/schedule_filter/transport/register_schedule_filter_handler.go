package transport

import (
	"github.com/gofiber/fiber/v2"
)

type ScheduleFilterHandlerRegister struct {
	Router fiber.Router
}

func RegisterScheduleFilterHandler(router fiber.Router) {
	scheduleFilterHandler := &ScheduleFilterHandlerRegister{
		Router: router,
	}

	// Register all endpoints here
	router.Get("/schedule", scheduleFilterHandler.ScheduleFilter)
}
