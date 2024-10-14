package comment

import (
	"api/dms"
	"encoding/json"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
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
