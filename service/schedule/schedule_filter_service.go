package schedule

import (
	"api/dms"
	"bytes"
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

	// Get user_email_id by userId
	userId := c.Locals("userid").(string)
	resp, err := dms.CallAPI("GET", "/user_email/user/"+userId, nil, nil, map[string]string{"status": ""}, 120)
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
	var userResponse []models.TwUserEmail
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, errors.New("could not unmarshal response body: " + err.Error())
	}
	// parse userResponse to get list of user_email_id
	user_email_id := make([]string, len(userResponse))
	for i, user := range userResponse {
		user_email_id[i] = strconv.Itoa(user.ID)
	}

	// Get workspace IDs for the user
	resp, err = dms.CallAPI("POST", "/workspace_user/user_email_id", user_email_id, nil, nil, 120)
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
	// có chứa role rồi
	var workspaceResponse []models.TwWorkspaceUser
	err = json.Unmarshal(body, &workspaceResponse)
	if err != nil {
		return nil, errors.New("could not unmarshal response body: " + err.Error())
	}
	// get wsp_user_id from workspaceResponse
	wspUserId := make([]string, len(workspaceResponse))
	for i, wsp := range workspaceResponse {
		wspUserId[i] = strconv.Itoa(wsp.ID)
	}
	// Validate workspace IDs
	wspIds := strings.Split(wspId, ",")
	for i, id := range wspIds {
		wspIds[i] = strings.TrimSpace(id)
	}
	var missingIds []string
	for _, id := range wspIds {
		if !contains(workspaceResponse, id) {
			missingIds = append(missingIds, id)
		}
	}
	if len(missingIds) > 0 {
		return nil, errors.New("some workspace IDs do not belong to the current user: " + strings.Join(missingIds, ", "))
	}

	userRoles := map[int]string{} // Map workspace_id -> role
	for _, wsp := range workspaceResponse {
		userRoles[wsp.WorkspaceId] = wsp.Role
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
	} else if isDeleted == "" {
		queryParams["is_deleted"] = "false"
	}
	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		queryParams["assigned_to"] = assignedTo
	}

	// Call the filter API
	resp, err = dms.CallAPI("GET", "/schedule/schedules/filter", nil, nil, queryParams, 120)
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
	var schedules []models.TwSchedule
	err = json.Unmarshal(body, &schedules)
	if err != nil {
		return nil, errors.New("could not unmarshal response body: " + err.Error())
	}

	// Filter schedules with visibility "private"
	filteredSchedules := []models.TwSchedule{}
	for _, schedule := range schedules {
		if schedule.Visibility == "private" {
			role := userRoles[schedule.WorkspaceId]
			if role == "admin" || role == "owner" {
				// Admin/Owner can view all private schedules
				filteredSchedules = append(filteredSchedules, schedule)
			} else {
				// Check if user is a participant in the schedule
				queryParams := map[string]string{
					"workspace_user_id": strings.Join(wspUserId, ","),
				}
				resp, err := dms.CallAPI("GET", "/schedule_participant/"+strconv.Itoa(schedule.ID)+"/participants", nil, nil, queryParams, 120)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, err
				}
				if resp.StatusCode == fiber.StatusNotFound {
					continue
				}
				if resp.StatusCode != fiber.StatusOK {
					return nil, errors.New("error from external service: " + string(body))
				}
				var participants []models.TwScheduleParticipant
				err = json.Unmarshal(body, &participants)
				if err != nil {
					return nil, errors.New("could not unmarshal response body: " + err.Error())
				}
				for _, participant := range participants {
					participantID := strconv.Itoa(participant.WorkspaceUserId)
					for _, wspID := range wspUserId {
						if participantID == wspID {
							filteredSchedules = append(filteredSchedules, schedule)
							break
						}
					}
				}
			}
		} else {
			// Public schedules are included
			filteredSchedules = append(filteredSchedules, schedule)
		}
	}

	// Return filtered schedules
	responseBody, err := json.Marshal(filteredSchedules)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: fiber.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil
}

// Helper function to check if a slice contains an item
func contains(slice []models.TwWorkspaceUser, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(strconv.Itoa(s.WorkspaceId)) == strings.TrimSpace(item) {
			return true
		}
	}
	return false
}
