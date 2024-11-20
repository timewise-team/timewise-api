package transport

import (
	"github.com/gofiber/fiber/v2"
)

type NotificationHandlerRegister struct {
	Router  fiber.Router
	Handler *NotificationHandler
}

func RegisterNotificationHandler(router fiber.Router) {
	handler := NewNotificationHandler()
	notificationHandler := &NotificationHandlerRegister{
		Handler: handler,
	}

	router.Get("", notificationHandler.Handler.GetNotifications)
}
