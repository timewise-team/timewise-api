package schedule

import (
	"api/dms"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type ScheduleFilterService struct {
}

func NewScheduleFilterService() *ScheduleFilterService {
	return &ScheduleFilterService{}
}
func (s *ScheduleFilterService) ScheduleFilter(c *fiber.Ctx) (*http.Response, error) {
	queryParams := map[string]string{}

	workspaceID := c.Query("workspace_id")
	if workspaceID != "" {
		queryParams["workspace_id"] = workspaceID
	}

	boardColumnID := c.Query("board_column_id")
	if boardColumnID != "" {
		queryParams["board_column_id"] = boardColumnID
	}

	title := c.Query("title")
	if title != "" {
		queryParams["title"] = title
	}

	startTime := c.Query("start_time")
	if startTime != "" {
		queryParams["start_time"] = startTime
	}

	endTime := c.Query("end_time")
	if endTime != "" {
		queryParams["end_time"] = endTime
	}

	location := c.Query("location")
	if location != "" {
		queryParams["location"] = location
	}

	createdBy := c.Query("created_by")
	if createdBy != "" {
		queryParams["created_by"] = createdBy
	}

	status := c.Query("status")
	if status != "" {
		queryParams["status"] = status
	}

	isDeleted := c.Query("is_deleted")
	if isDeleted != "" {
		queryParams["is_deleted"] = isDeleted
	}

	assignedTo := c.Query("assigned_to")
	if assignedTo != "" {
		queryParams["assigned_to"] = assignedTo
	}

	resp, err := dms.CallAPI("GET", "/schedule/schedules/filter", nil, nil, queryParams, 120)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
