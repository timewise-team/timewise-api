package transport

import (
	"github.com/gofiber/fiber/v2"
)

type DocumentHandlerRegister struct {
	Router  fiber.Router
	Handler *DocumentHandler
}

func RegisterDocumentHandler(router fiber.Router) {
	handler := NewDocumentHandler()
	scheduleHandler := &DocumentHandlerRegister{
		Handler: handler,
	}

	// Register all endpoints here
	router.Get("schedule/:scheduleId",
		scheduleHandler.Handler.GetDocumentByScheduleID)
}
