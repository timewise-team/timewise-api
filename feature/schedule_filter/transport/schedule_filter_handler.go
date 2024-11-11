package transport

import (
	"api/dms"
	"api/service/schedule"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	dtos "github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io"
	"strconv"
)

type ScheduleFilterHandler struct {
	service schedule.ScheduleFilterService
}

func NewScheduleFilterHandler() *ScheduleFilterHandler {
	service := schedule.NewScheduleFilterService()
	return &ScheduleFilterHandler{
		service: *service,
	}
}

// ScheduleFilter godoc
// @Summary Get schedules by filter
// @Description Retrieve a list of schedules based on specified filter parameters.
// @Tags schedule
// @Accept json
// @Produce json
// @Security bearerToken
// @Param workspace_id query string false "Filter by Workspace ID"
// @Param board_column_id query int false "Filter by Board Column ID"
// @Param title query string false "Filter by Title (searches with LIKE)"
// @Param start_time query string false "Filter by Start Time (ISO8601 format)"
// @Param end_time query string false "Filter by End Time (ISO8601 format)"
// @Param location query string false "Filter by Location (searches with LIKE)"
// @Param created_by query int false "Filter by User ID of the creator"
// @Param status query string false "Filter by Status"
// @Param is_deleted query bool false "Filter by Deleted Schedules"
// @Param assigned_to query int false "Filter by User ID assigned to the schedule"
// @Success 200 {array} dtos.TwScheduleResponse "List of filtered schedules"
// @Failure 400 {object} fiber.Error "Invalid query parameters"
// @Failure 500 {object} fiber.Error "Internal Server Error"
// @Router /api/v1/schedule/schedule [get]
func (h *ScheduleFilterHandler) ScheduleFilter(c *fiber.Ctx) error {
	// check wsp id
	wspId := c.Query("workspace_id")
	if wspId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Workspace ID is required",
		})
	}
	// get user_email_id by email
	email := c.Locals("email").(string)
	resp, err := dms.CallAPI("GET", "user_email/email"+email, nil, nil, nil, 120)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch user_email_id from service",
			"details": err.Error(),
		})
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response body",
		})
	}
	if resp.StatusCode != fiber.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error":   "Error from external service",
			"details": string(body),
		})
	}
	var userResponse models.TwUserEmail
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Could not unmarshal response body",
			"details": err.Error(),
		})
	}
	// get wsp_id belong to current account in workspace_user
	resp, err = dms.CallAPI("GET", "workspace_user/user_email_id/"+strconv.Itoa(userResponse.ID), nil, nil, nil, 120)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch workspace_id from service",
			"details": err.Error(),
		})
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response body",
		})
	}
	if resp.StatusCode != fiber.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error":   "Error from external service",
			"details": string(body),
		})
	}
	var workspaceUserResponse []models.TwWorkspaceUser
	err = json.Unmarshal(body, &workspaceUserResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Could not unmarshal response body",
			"details": err.Error(),
		})
	}
	// check wsp_id belong to current account
	var check bool
	for _, workspaceUser := range workspaceUserResponse {
		if workspaceUser.WorkspaceID == wspId {
			check = true
			break
		}
	}
	resp, err = h.service.ScheduleFilter(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch schedules from service",
			"details": err.Error(),
		})
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response body",
		})
	}

	if resp.StatusCode >= 400 {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error":   "Error from external service",
			"details": string(body),
		})
	}

	var scheduleResponse []dtos.TwScheduleResponse

	err = json.Unmarshal(body, &scheduleResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Could not unmarshal response body",
			"details": err.Error(),
		})
	}

	return c.Status(resp.StatusCode).JSON(scheduleResponse)
}
