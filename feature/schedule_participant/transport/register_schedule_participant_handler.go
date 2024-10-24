package transport

import (
	"github.com/gofiber/fiber/v2"
)

type ScheduleParticipantHandlerRegister struct {
	Router  fiber.Router
	Handler *ScheduleParticipantHandler
}

func RegisterScheduleParticipantHandler(router fiber.Router) {
	handler := NewScheduleParticipantHandler()
	scheduleParticipantHandler := &ScheduleParticipantHandlerRegister{
		Handler: handler,
	}

	// Register all endpoints here
	router.Get("schedule/:scheduleId",
		scheduleParticipantHandler.Handler.GetScheduleParticipantByScheduleID)
}
