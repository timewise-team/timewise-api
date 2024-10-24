package transport

import (
	"api/service/comment"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CommentHandler struct {
	service comment.CommentService
}

func NewCommentHandler() *CommentHandler {
	service := comment.NewCommentService()
	return &CommentHandler{
		service: *service,
	}
}

// getCommentsBySchedule godoc
// @Summary Get comments by schedule
// @Description Get comments by schedule
// @Tags comment
// @Accept json
// @Produce json
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {array} comment_dtos.TwCommentResponse
// @Router /api/v1/comment/schedule/{schedule_id} [get]
func (h *CommentHandler) GetCommentByScheduleID(c *fiber.Ctx) error {
	scheduleIDStr := c.Params("scheduleID")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid schedule ID")
	}
	comment, err := h.service.GetCommentsByScheduleID(scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(comment)
}
