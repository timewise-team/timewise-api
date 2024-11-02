package transport

import (
	"api/service/comment"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/comment_dtos"
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

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a new comment
// @Tags comment
// @Accept json
// @Produce json
// @Param schedule body comment_dtos.CommentRequestDTO true "Comment"
// @Success 201 {object} comment_dtos.CommentResponseDTO
// @Router /api/v1/comment [post]
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	var CreateCommentDto comment_dtos.CommentRequestDTO
	if err := c.BodyParser(&CreateCommentDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	result, err := h.service.CreateComment(c, CreateCommentDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(result)
}

// UpdateComment godoc
// @Summary Update an existing comment
// @Description Update an existing comment
// @Tags comment
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Param schedule body comment_dtos.CommentRequestDTO true "Comment"
// @Success 200 {object} comment_dtos.CommentResponseDTO
// @Router /api/v1/comment/{comment_id} [put]
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	var UpdateCommentDto comment_dtos.CommentRequestDTO
	if err := c.BodyParser(&UpdateCommentDto); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	result, err := h.service.UpdateComment(c, UpdateCommentDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(result)
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Delete a comment
// @Tags comment
// @Accept json
// @Produce json
// @Param comment_id path int true "Comment ID"
// @Success 204 "No Content"
// @Router /api/v1/comment/{comment_id} [delete]
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	result, err := h.service.DeleteComment(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(result)
}
