package transport

import "github.com/gofiber/fiber/v2"

type UserEmailsHandler struct {
	Router fiber.Router
}

func RegisterUserEmailsHandler(router fiber.Router) {
	userEmailsHandler := UserEmailsHandler{
		Router: router,
	}

	// Register all endpoints here
	router.Post("/", userEmailsHandler.createNewUserEmail)
	//router.Post("/logout", authHandler.logout)
	//router.Post("/refresh", authHandler.refresh)
	//router.Post("/forgot-password", authHandler.forgotPassword)
	//router.Post("/reset-password", authHandler.resetPassword)
}
