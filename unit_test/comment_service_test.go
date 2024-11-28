package unit_test_test

import (
	"api/service/comment"
	"api/unit_test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/comment_dtos"
	"github.com/timewise-team/timewise-models/models"
	"testing"
)

func TestFunc27_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	commenter := "hihi"
	content := "hello"
	scheduleId := -1

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := comment_dtos.CommentRequestDTO{
		Commenter:  &commenter,
		Content:    &content,
		ScheduleId: &scheduleId,
	}

	_, err := service.CreateComment(&workspaceUser, request)

	assert.Equal(t, "schedule id must be greater than zero", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc27_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	commenter := "hihi"
	content := "hello"
	scheduleId := 0

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := comment_dtos.CommentRequestDTO{
		Commenter:  &commenter,
		Content:    &content,
		ScheduleId: &scheduleId,
	}

	_, err := service.CreateComment(&workspaceUser, request)

	assert.Equal(t, "schedule id must be greater than zero", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc27_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	commenter := "hihi"
	content := "hello"
	scheduleId := 102

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := comment_dtos.CommentRequestDTO{
		Commenter:  &commenter,
		Content:    &content,
		ScheduleId: &scheduleId,
	}

	createdComment, err := service.CreateComment(&workspaceUser, request)

	assert.NoError(t, err)
	assert.Equal(t, "hihi", createdComment.Commenter)
	assert.Equal(t, "hello", createdComment.Content)
	assert.Equal(t, 102, createdComment.ScheduleId)
	mockDMS.AssertExpectations(t)
}

func TestFunc27_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	commenter := "hihi"
	content := ""
	scheduleId := 102

	workspaceUser := models.TwWorkspaceUser{
		ID: 2,
	}

	request := comment_dtos.CommentRequestDTO{
		Commenter:  &commenter,
		Content:    &content,
		ScheduleId: &scheduleId,
	}

	_, err := service.CreateComment(&workspaceUser, request)

	assert.Equal(t, "content is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc28_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	content := "hello"

	request := comment_dtos.CommentRequestDTO{
		Content: &content,
	}

	_, err := service.UpdateComment("-1", request)

	assert.Equal(t, "comment id must be greater than zero", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc28_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	content := "hello"

	request := comment_dtos.CommentRequestDTO{
		Content: &content,
	}

	_, err := service.UpdateComment("0", request)

	assert.Equal(t, "comment id must be greater than zero", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc28_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	content := "hi my name is Thuan"

	request := comment_dtos.CommentRequestDTO{
		Content: &content,
	}

	updatedComment, _ := service.UpdateComment("18", request)

	assert.Equal(t, "hi my name is Thuan", updatedComment.Content)
	mockDMS.AssertExpectations(t)
}

func TestFunc28_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	content := ""

	request := comment_dtos.CommentRequestDTO{
		Content: &content,
	}

	_, err := service.UpdateComment("18", request)

	assert.Equal(t, "content is required", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc29_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	_, err := service.GetCommentsByScheduleID(-1)

	assert.Equal(t, "schedule id must be greater than zero", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc29_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	_, err := service.GetCommentsByScheduleID(0)

	assert.Equal(t, "schedule id must be greater than zero", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc29_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	listComment, _ := service.GetCommentsByScheduleID(5)

	assert.Equal(t, 7, len(listComment))
	mockDMS.AssertExpectations(t)
}

func TestFunc30_UTCID01(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	_, err := service.DeleteComment("-1")

	assert.Equal(t, "comment id must be greater than zero", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc30_UTCID02(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	_, err := service.DeleteComment("0")

	assert.Equal(t, "comment id must be greater than zero", err.Error())
	mockDMS.AssertExpectations(t)
}

func TestFunc30_UTCID03(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	deletedComment, _ := service.DeleteComment("18")

	assert.Equal(t, true, deletedComment.IsDeleted)
	mockDMS.AssertExpectations(t)
}

func TestFunc30_UTCID04(t *testing.T) {
	utils.InitConfig()
	mockDMS := new(MockDMSClient)
	service := comment.NewCommentService()

	_, err := service.DeleteComment("100")

	assert.Equal(t, "comment not found", err.Error())
	mockDMS.AssertExpectations(t)
}
