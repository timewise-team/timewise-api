package transport

import (
	"api/middleware"
	"github.com/gofiber/fiber/v2"
)

type CommentHandlerRegister struct {
	Router  fiber.Router
	Handler *CommentHandler
}

func RegisterCommentHandler(router fiber.Router) {
	handler := NewCommentHandler()
	commentHandler := &CommentHandlerRegister{
		Handler: handler,
	}

	// Register all endpoints here
	router.Get("schedule/:scheduleId",
		commentHandler.Handler.GetCommentByScheduleID)
	router.Post("/",
		middleware.CheckWorkspaceRole([]string{"owner", "admin", "member"}),
		middleware.CheckScheduleStatus([]string{"creator", "assign to"}),
		commentHandler.Handler.CreateComment)
	router.Put("/:id", commentHandler.Handler.UpdateComment)
	router.Delete("/:id", commentHandler.Handler.DeleteComment)
}
