package transport

import (
	"github.com/gofiber/fiber/v2"
)

type BoardColumnsHandler struct {
	Router fiber.Router
}

func RegisterAuthHandler(router fiber.Router) {
	boardColumnsHandler := BoardColumnsHandler{
		Router: router,
	}

	// Register all endpoints here
	router.Post("/workspace/:workspace_id/board_columns", boardColumnsHandler.getBoardColumnsByWorkspace)
	//router.Post("/logout", authHandler.logout)
	//router.Post("/refresh", authHandler.refresh)
	//router.Post("/forgot-password", authHandler.forgotPassword)
	//router.Post("/reset-password", authHandler.resetPassword)
}
