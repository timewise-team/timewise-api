package schedule

import (
	"api/dms"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

type ScheduleFilterService struct {
}

func NewScheduleFilterService() *ScheduleFilterService {
	return &ScheduleFilterService{}
}
func (s *ScheduleFilterService) ScheduleFilter(c *fiber.Ctx) (*http.Response, error) {
	queryParams := map[string]string{}

	workspaceIDs := c.Context().QueryArgs().PeekMulti("workspace_id")
	if len(workspaceIDs) > 0 {
		var workspaceIDStrings []string
		for _, id := range workspaceIDs {
			workspaceIDStrings = append(workspaceIDStrings, string(id)) // Convert each []byte to string
		}

		// Join the workspace_ids into a single string with commas as delimiter
		queryParams["workspace_id"] = strings.Join(workspaceIDStrings, ",")
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
