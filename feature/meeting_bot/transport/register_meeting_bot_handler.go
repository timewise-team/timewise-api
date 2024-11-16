package transport

import (
	"github.com/gofiber/fiber/v2"
)

type MeetingBotHandlerRegister struct {
	Router  fiber.Router
	Handler *MeetingBotHandler
}

func RegisterMeetingBotHandler(router fiber.Router) {
	handler, _ := NewMeetingBotHandler()
	reminderHandler := &MeetingBotHandlerRegister{
		Handler: handler,
	}

	router.Post("/start", reminderHandler.StartMeeting)
}
