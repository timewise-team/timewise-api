package transport

import "github.com/gofiber/fiber/v2"

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
}
