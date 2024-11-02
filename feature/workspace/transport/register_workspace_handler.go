package transport

import (
	"api/middleware"
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
	router.Delete("/delete-workspace", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), workspaceHandler.deleteWorkspace)
	router.Get("/get-workspace-by-id/:workspace_id", workspaceHandler.getWorkspaceById)
	router.Put("/update-workspace", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), workspaceHandler.updateWorkspace)
	router.Get("/filter-workspaces", workspaceHandler.filterWorkspaces)
	//router.Post("/logout", authHandler.logout)
	//router.Post("/refresh", authHandler.refresh)
	//router.Post("/forgot-password", authHandler.forgotPassword)
	//router.Post("/reset-password", authHandler.resetPassword)
}
