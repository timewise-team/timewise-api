package transport

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

type NotificationRegister struct {
	Router  fiber.Router
	Handler *NotificationHandler
}

// RegisterNotificationHandler initializes and registers the notification handler routes
func RegisterNotificationHandler(router fiber.Router) {
	handler, err := NewNotificationHandler()
	if err != nil {
		log.Println("Error initializing NotificationHandler:", err)
		return
	}

	notificationHandler := &NotificationRegister{
		Handler: handler,
	}

	// Register all notification endpoints
	router.Post("/push", notificationHandler.Handler.PushNotifications)
}
