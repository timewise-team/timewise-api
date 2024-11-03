package comment

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/comment_dtos"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
	"time"
)

type CommentService struct {
}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (h *CommentService) GetCommentsBySchedule(scheduleId int) ([]models.TwComment, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/comment/schedule/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var comments []models.TwComment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (h *CommentService) GetCommentsByScheduleID(scheduleId int) ([]comment_dtos.TwCommentResponse, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/comment/schedule_id/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var comments []comment_dtos.TwCommentResponse
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentService) CreateComment(c *fiber.Ctx, CommentRequestDto comment_dtos.CommentRequestDTO) (*comment_dtos.CommentResponseDTO, error) {

	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return nil, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve schedule participant",
		})
	}

	newComment := models.TwComment{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ScheduleId:      *CommentRequestDto.ScheduleId,
		WorkspaceUserId: workspaceUser.ID,
		Commenter:       *CommentRequestDto.Commenter,
		Content:         *CommentRequestDto.Content,
	}

	resp, err := dms.CallAPI(
		"POST",
		"/comment",
		newComment,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var comment models.TwComment
	if err := json.NewDecoder(resp.Body).Decode(&comment); err != nil {
		return nil, err
	}

	newCommentResponse := comment_dtos.CommentResponseDTO{
		ID:              comment.ID,
		CreatedAt:       &comment.CreatedAt,
		UpdatedAt:       &comment.UpdatedAt,
		ScheduleId:      comment.ScheduleId,
		WorkspaceUserId: comment.WorkspaceUserId,
		Commenter:       comment.Commenter,
		Content:         comment.Content,
	}

	return &newCommentResponse, nil
}

func (s *CommentService) UpdateComment(c *fiber.Ctx, CommentRequestDto comment_dtos.CommentRequestDTO) (*comment_dtos.CommentResponseDTO, error) {
	commentId := c.Params("id")
	resp1, err1 := dms.CallAPI(
		"GET",
		"/comment/"+commentId,
		nil,
		nil,
		nil,
		120,
	)
	if err1 != nil {
		return nil, err1
	}
	defer resp1.Body.Close()

	if resp1.StatusCode == fiber.StatusNotFound {
		return nil, errors.New("comment not found")
	}

	var comment models.TwComment
	if err := json.NewDecoder(resp1.Body).Decode(&comment); err != nil {
		return nil, err
	}

	newComment := models.TwComment{
		ID:              comment.ID,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       time.Now(),
		ScheduleId:      comment.ScheduleId,
		WorkspaceUserId: comment.WorkspaceUserId,
		Commenter:       comment.Commenter,
		Content:         *CommentRequestDto.Content,
	}

	resp, err := dms.CallAPI(
		"PUT",
		"/comment/"+commentId,
		newComment,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateComment models.TwComment
	if err := json.NewDecoder(resp.Body).Decode(&updateComment); err != nil {
		return nil, err
	}

	newCommentResponse := comment_dtos.CommentResponseDTO{
		ID:              updateComment.ID,
		CreatedAt:       &updateComment.CreatedAt,
		UpdatedAt:       &updateComment.UpdatedAt,
		ScheduleId:      updateComment.ScheduleId,
		WorkspaceUserId: updateComment.WorkspaceUserId,
		Commenter:       updateComment.Commenter,
		Content:         updateComment.Content,
	}
	return &newCommentResponse, nil
}

func (s *CommentService) DeleteComment(c *fiber.Ctx) (*comment_dtos.CommentResponseDTO, error) {
	commentId := c.Params("id")
	resp1, err1 := dms.CallAPI(
		"GET",
		"/comment/"+commentId,
		nil,
		nil,
		nil,
		120,
	)
	if err1 != nil {
		return nil, err1
	}
	defer resp1.Body.Close()

	if resp1.StatusCode == fiber.StatusNotFound {
		return nil, errors.New("comment not found")
	}

	var comment models.TwComment
	if err := json.NewDecoder(resp1.Body).Decode(&comment); err != nil {
		return nil, err
	}

	now := time.Now()

	deleteComment := models.TwComment{
		ID:              comment.ID,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       time.Now(),
		DeletedAt:       &now,
		ScheduleId:      comment.ScheduleId,
		WorkspaceUserId: comment.WorkspaceUserId,
		Commenter:       comment.Commenter,
		Content:         comment.Content,
		IsDeleted:       true,
	}

	resp, err := dms.CallAPI(
		"DELETE",
		"/comment/"+commentId,
		deleteComment,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateComment models.TwComment
	if err := json.NewDecoder(resp.Body).Decode(&updateComment); err != nil {
		return nil, err
	}

	newCommentResponse := comment_dtos.CommentResponseDTO{
		ID:              updateComment.ID,
		CreatedAt:       &updateComment.CreatedAt,
		UpdatedAt:       &updateComment.UpdatedAt,
		DeletedAt:       updateComment.DeletedAt,
		ScheduleId:      updateComment.ScheduleId,
		WorkspaceUserId: updateComment.WorkspaceUserId,
		Commenter:       updateComment.Commenter,
		Content:         updateComment.Content,
		IsDeleted:       updateComment.IsDeleted,
	}
	return &newCommentResponse, nil
}
