package transport

import (
	"github.com/gofiber/fiber/v2"
)

type WorkspaceHandler struct {
	Router fiber.Router
}

func RegisterWorkspaceHandler(router fiber.Router) {
	workspaceHandler := WorkspaceHandler{
		Router: router,
	}

	// Register all endpoints here
	router.Get("/get-workspaces-by-email/:email?", workspaceHandler.getWorkspacesByEmail)
	router.Post("/create-workspace", workspaceHandler.createWorkspace)
	//router.Post("/logout", authHandler.logout)
	//router.Post("/refresh", authHandler.refresh)
	//router.Post("/forgot-password", authHandler.forgotPassword)
	//router.Post("/reset-password", authHandler.resetPassword)
}
