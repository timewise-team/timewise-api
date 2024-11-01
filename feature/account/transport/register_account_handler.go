package transport

import "github.com/gofiber/fiber/v2"

type AccountHandlerRegister struct {
	Router  fiber.Router
	Handler *AccountHandler
}

func RegisterAccountHandler(router fiber.Router) {
	handler := NewAccountHandler()
	accountHandler := &AccountHandlerRegister{
		Handler: handler,
	}

	// Register all endpoints here
	router.Get("/user", accountHandler.Handler.getUserInfo)
	router.Get("/user/emails", accountHandler.Handler.getLinkedUserEmails)
	router.Patch("/user", accountHandler.Handler.updateUserInfo)
	router.Post("/user/emails", accountHandler.Handler.linkAnEmail)
	router.Post("/user/emails/unlink", accountHandler.Handler.unlinkAnEmail)
}
