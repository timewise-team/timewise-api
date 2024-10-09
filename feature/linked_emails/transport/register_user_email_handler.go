package transport

import (
	"github.com/gofiber/fiber/v2"
)

type LinkedEmailsHandler struct {
	Router fiber.Router
}

func RegisterLinkedEmailsHandler(router fiber.Router) {
	linkedEmailsHandler := LinkedEmailsHandler{
		Router: router,
	}

	// Register all endpoints here
	router.Get("/get-linked-email", linkedEmailsHandler.getLinkedUserEmail)

}
