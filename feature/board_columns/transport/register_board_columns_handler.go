package transport

import (
	"api/middleware"
	"github.com/gofiber/fiber/v2"
)

type BoardColumnsHandler struct {
	Router fiber.Router
}

func RegisterBoardColumnsHandler(router fiber.Router) {
	boardColumnsHandler := BoardColumnsHandler{
		Router: router,
	}

	// Register all endpoints here
	router.Get("/workspace/:workspace_id", middleware.CheckWorkspaceRole([]string{"owner", "admin", "member", "guest"}), boardColumnsHandler.getBoardColumnsByWorkspace)
	router.Get("/workspace_id/:workspace_id", boardColumnsHandler.getBoardColumnsByWorkspaceId)
	router.Post("", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), boardColumnsHandler.createBoardColumn)
	router.Delete("/:board_column_id", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), boardColumnsHandler.deleteBoardColumn)
	router.Put("/:board_column_id", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), boardColumnsHandler.updateBoardColumn)
	router.Put("/update_position/:board_column_id", middleware.CheckWorkspaceRole([]string{"owner", "admin"}), boardColumnsHandler.updatePosition)
	//router.Post("/logout", authHandler.logout)
	//router.Post("/refresh", authHandler.refresh)
	//router.Post("/forgot-password", authHandler.forgotPassword)
	//router.Post("/reset-password", authHandler.resetPassword)
}
