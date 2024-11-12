package schedule

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ScheduleFilterService struct{}

func NewScheduleFilterService() *ScheduleFilterService {
	return &ScheduleFilterService{}
}

func (s *ScheduleFilterService) ScheduleFilter(c *fiber.Ctx) (*http.Response, error) {
	// Check for workspace_id
	wspId := c.Query("workspace_id")
	if wspId == "" {
		return nil, errors.New("Workspace ID is required")
	}

	// Get user_email_id by email
	email := c.Locals("email").(string)
	resp, err := dms.CallAPI("GET", "/user_email/email/"+email, nil, nil, nil, 120)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != fiber.StatusOK {
		return nil, errors.New("error from external service: " + string(body))
	}
	var userResponse models.TwUserEmail
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, errors.New("could not unmarshal response body: " + err.Error())
	}

	// Get workspace IDs for the user
	resp, err = dms.CallAPI("GET", "/workspace_user/user_email_id/"+strconv.Itoa(userResponse.ID), nil, nil, nil, 120)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != fiber.StatusOK {
		return nil, errors.New("error from external service: " + string(body))
	}
	var workspaceIdResponse []string
	err = json.Unmarshal(body, &workspaceIdResponse)
	if err != nil {
		return nil, errors.New("could not unmarshal response body: " + err.Error())
	}

	// Validate workspace IDs
	wspIds := strings.Split(wspId, ",")
	for i, id := range wspIds {
		wspIds[i] = strings.TrimSpace(id)
	}
	var missingIds []string
	for _, id := range wspIds {
		if !contains(workspaceIdResponse, id) {
			missingIds = append(missingIds, id)
		}
	}
	if len(missingIds) > 0 {
		return nil, errors.New("some workspace IDs do not belong to the current user: " + strings.Join(missingIds, ", "))
	}

	// Construct queryParams
	queryParams := map[string]string{
		"workspace_id": strings.Join(wspIds, ","),
	}
	if title := c.Query("title"); title != "" {
		queryParams["title"] = title
	}
	if startTime := c.Query("start_time"); startTime != "" {
		queryParams["start_time"] = startTime
	}
	if endTime := c.Query("end_time"); endTime != "" {
		queryParams["end_time"] = endTime
	}
	if location := c.Query("location"); location != "" {
		queryParams["location"] = location
	}
	if createdBy := c.Query("created_by"); createdBy != "" {
		queryParams["created_by"] = createdBy
	}
	if status := c.Query("status"); status != "" {
		queryParams["status"] = status
	}
	if isDeleted := c.Query("is_deleted"); isDeleted != "" {
		queryParams["is_deleted"] = isDeleted
	}
	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		queryParams["assigned_to"] = assignedTo
	}

	// Call the filter API
	resp, err = dms.CallAPI("GET", "/schedule/schedules/filter", nil, nil, queryParams, 120)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Helper function to check if a slice contains an item
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == strings.TrimSpace(item) {
			return true
		}
	}
	return false
}
