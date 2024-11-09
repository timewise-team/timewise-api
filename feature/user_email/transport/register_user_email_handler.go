package transport

import (
	"github.com/gofiber/fiber/v2"
)

type UserEmailHandler struct {
	Router fiber.Router
}

func RegisterUserEmailHandler(router fiber.Router) {
	userEmailHandler := UserEmailHandler{
		Router: router,
	}

	router.Get("/search-user_email/:query?", userEmailHandler.searchUserEmail)
	router.Get("/list_approve/:scheduleId", userEmailHandler.listApproveUserEmailHandler)
}
