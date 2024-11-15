package transport

import (
	"github.com/gofiber/fiber/v2"
)

type NotificationSettingHandlerRegister struct {
	Router  fiber.Router
	Handler *NotificationSettingHandler
}

func RegisterNotificationSettingHandler(router fiber.Router) {
	handler := NewNotificationSettingHandler()
	notificationSettingHandler := &NotificationSettingHandlerRegister{
		Handler: handler,
	}

	router.Get("", notificationSettingHandler.Handler.GetNotificationSettingByUserId)
	router.Put("", notificationSettingHandler.Handler.UpdateNotificationSetting)
}
