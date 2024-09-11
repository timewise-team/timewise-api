package authentication

import (
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Router fiber.Router
}

func RegisterAuthHandler(router fiber.Router) {
	authHandler := AuthHandler{
		Router: router,
	}

	// Register all endpoints here
	router.Post("/login", authHandler.login)
	//router.Post("/register", authHandler.register)
	//router.Post("/logout", authHandler.logout)
	//router.Post("/refresh", authHandler.refresh)
	//router.Post("/forgot-password", authHandler.forgotPassword)
	//router.Post("/reset-password", authHandler.resetPassword)
}
