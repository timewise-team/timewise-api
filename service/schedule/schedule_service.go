package schedule

import (
	"api/dms"
	"api/notification"
	"api/service/reminder"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ScheduleService struct {
}

func NewScheduleService() *ScheduleService {
	return &ScheduleService{}
}

func (s *ScheduleService) CreateSchedule(workspaceUser *models.TwWorkspaceUser, CreateScheduleDto core_dtos.TwCreateScheduleRequest) (*core_dtos.TwCreateShecduleResponse, int, error) {

	if CreateScheduleDto.BoardColumnID == nil || *CreateScheduleDto.BoardColumnID == 0 {
		return nil, http.StatusBadRequest, errors.New("Invalid board column id")
	}
	if CreateScheduleDto.WorkspaceID == nil || *CreateScheduleDto.WorkspaceID == 0 {
		return nil, http.StatusBadRequest, errors.New("Invalid workspace id")
	}
	if *CreateScheduleDto.Title == "" {
		return nil, http.StatusBadRequest, errors.New("Invalid title")
	}
	newSchedule := core_dtos.TwCreateScheduleRequest{
		WorkspaceID:     CreateScheduleDto.WorkspaceID,
		BoardColumnID:   CreateScheduleDto.BoardColumnID,
		WorkspaceUserID: &workspaceUser.ID,
		Title:           CreateScheduleDto.Title,
		Description:     CreateScheduleDto.Description,
		StartTime:       CreateScheduleDto.StartTime,
		EndTime:         CreateScheduleDto.EndTime,
	}

	resp1, err := dms.CallAPI(
		"POST",
		"/schedule",
		newSchedule,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	defer resp1.Body.Close()

	var createdSchedule core_dtos.TwCreateShecduleResponse
	if err := json.NewDecoder(resp1.Body).Decode(&createdSchedule); err != nil {
		return nil, 500, fmt.Errorf("error decoding schedule response: %v", err)
	}

	scheduleDetail, err := s.GetScheduleDetailByID(strconv.Itoa(createdSchedule.ID))
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}
	if scheduleDetail.StartTime != nil {
		startTime := *scheduleDetail.StartTime
		err1 := reminder.NewReminderService().CreateReminderAllParticipantWhenCreateSchedule(createdSchedule.ID, startTime, workspaceUser, 0)
		if err1 != nil {
			return nil, fiber.StatusInternalServerError, err
		}
	}

	// send notification
	notificationDto := models.TwNotifications{
		Title:       "New Schedule created",
		Description: fmt.Sprintf("You have created new schedule %s", scheduleDetail.Title),
		Link:        fmt.Sprintf("/organization/%d?schedule_id=%d", *CreateScheduleDto.WorkspaceID, createdSchedule.ID),
		UserEmailId: workspaceUser.UserEmailId,
		Type:        "schedule_created",
		IsSent:      true,
	}
	err = notification.PushNotifications(notificationDto)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return &createdSchedule, resp1.StatusCode, nil
}

func fetchSchedule(scheduleID string) (models.TwSchedule, error) {
	resp, err := dms.CallAPI("GET", "/schedule/"+scheduleID, nil, nil, nil, 120)
	if err != nil {
		return models.TwSchedule{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return models.TwSchedule{}, fmt.Errorf("GET /schedule/%s returned status %d: %s", scheduleID, resp.StatusCode, string(body))
	}

	var schedule models.TwSchedule
	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		return models.TwSchedule{}, fmt.Errorf("error decoding schedule response: %v", err)
	}

	return schedule, nil
}

func (s *ScheduleService) FetchScheduleParticipant(workspaceUserIdStr, scheduleID string) (models.TwScheduleParticipant, error) {
	resp, err := dms.CallAPI("GET", "/schedule_participant/workspace_user/"+workspaceUserIdStr+"/schedule/"+scheduleID, nil, nil, nil, 120)
	if err != nil {
		return models.TwScheduleParticipant{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return models.TwScheduleParticipant{}, fmt.Errorf("GET /schedule_participant/workspace_user/%s/schedule/%s returned status %d: %s", workspaceUserIdStr, scheduleID, resp.StatusCode, string(body))
	}

	var scheduleParticipant models.TwScheduleParticipant
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipant); err != nil {
		return models.TwScheduleParticipant{}, fmt.Errorf("error decoding schedule_participant response: %v", err)
	}

	return scheduleParticipant, nil
}

func applyUpdateFields(baseSchedule, updateSchedule models.TwSchedule, dto core_dtos.TwUpdateScheduleRequest) (models.TwSchedule, error) {
	// Validate Title
	if dto.Title != nil {
		if *dto.Title == "" {
			return updateSchedule, fmt.Errorf("title cannot be empty")
		}
		updateSchedule.Title = *dto.Title
	}

	if dto.StartTime != nil {
		parsedStartTime, err := time.Parse("2006-01-02 15:04:05.000", *dto.StartTime)
		if err != nil {
			return updateSchedule, fmt.Errorf("error parsing start time: %v", err)
		}
		updateSchedule.StartTime = &parsedStartTime
	}

	// Validate EndTime
	if dto.EndTime != nil {
		// Parse EndTime
		parsedEndTime, err := time.Parse("2006-01-02 15:04:05.000", *dto.EndTime)
		if err != nil {
			return updateSchedule, fmt.Errorf("error parsing end time: %v", err)
		}
		// Kiểm tra EndTime không được nhỏ hơn StartTime
		if parsedEndTime.Before(*updateSchedule.StartTime) || parsedEndTime.Equal(*updateSchedule.StartTime) {
			return updateSchedule, fmt.Errorf("Invalid Endtime")
		}
		updateSchedule.EndTime = &parsedEndTime
	}

	// Apply other fields
	if dto.Description != nil {
		updateSchedule.Description = *dto.Description
	}
	if dto.Location != nil {
		updateSchedule.Location = *dto.Location
	}
	if dto.Status != nil {
		updateSchedule.Status = *dto.Status
	}
	if dto.AllDay != nil {
		updateSchedule.AllDay = *dto.AllDay
	}
	if dto.Visibility != nil {
		updateSchedule.Visibility = *dto.Visibility
	}
	if dto.ExtraData != nil {
		updateSchedule.ExtraData = ""
	}
	if dto.RecurrencePattern != nil {
		updateSchedule.RecurrencePattern = *dto.RecurrencePattern
	}
	if dto.Priority != nil {
		updateSchedule.Priority = *dto.Priority
	}

	return updateSchedule, nil
}

func (s *ScheduleService) UpdateSchedule(
	scheduleID string,
	scheduleParticipant models.TwScheduleParticipant,
	workspaceUser *models.TwWorkspaceUser,
	UpdateScheduleDto core_dtos.TwUpdateScheduleRequest) (*models.TwSchedule, error) {

	if scheduleID == "" || scheduleID == "0" {
		return nil, fmt.Errorf("schedule id is required")
	}
	schedule, err := fetchSchedule(scheduleID)
	if err != nil {
		return nil, err
	}

	updateSchedule := schedule
	updateSchedule, err = applyUpdateFields(schedule, updateSchedule, UpdateScheduleDto)
	if err != nil {
		return nil, fmt.Errorf("Bad Request: %v", err)
	}

	if scheduleParticipant.Status == "assign to" {
		updateSchedule.Title = schedule.Title
		updateSchedule.Description = schedule.Description
		updateSchedule.Location = schedule.Location
		updateSchedule.Visibility = schedule.Visibility
		updateSchedule.VideoTranscript = schedule.VideoTranscript
		updateSchedule.ExtraData = ""
		updateSchedule.RecurrencePattern = schedule.RecurrencePattern
		updateSchedule.StartTime = schedule.StartTime
		updateSchedule.EndTime = schedule.EndTime
		updateSchedule.Priority = schedule.Priority

		if UpdateScheduleDto.Status != nil {
			updateSchedule.Status = *UpdateScheduleDto.Status
		}

	}

	resp, err := dms.CallAPI("PUT", "/schedule/"+scheduleID+"/workspace_user/"+strconv.Itoa(workspaceUser.ID), updateSchedule, nil, nil, 120)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("cannot read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var updatedSchedule models.TwSchedule
	if err := json.Unmarshal(body, &updatedSchedule); err != nil {
		return nil, errors.New("error parsing JSON")
	}

	return &updatedSchedule, nil
}

func (s *ScheduleService) UpdateSchedulePosition(scheduleId string, workspaceUser *models.TwWorkspaceUser, UpdateScheduleDto core_dtos.TwUpdateSchedulePosition) (*core_dtos.TwUpdateScheduleResponse, error) {

	if *UpdateScheduleDto.Position == 0 || *UpdateScheduleDto.Position == -1 {
		return nil, fmt.Errorf("invalid position")
	}
	if *UpdateScheduleDto.BoardColumnID == 0 || *UpdateScheduleDto.BoardColumnID == -1 {
		return nil, fmt.Errorf("invalid board column id")
	}
	resp, err := dms.CallAPI("PUT", "/schedule/position/"+scheduleId+"/workspace_user/"+strconv.Itoa(workspaceUser.ID), UpdateScheduleDto, nil, nil, 120)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("cannot read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var updatedSchedule core_dtos.TwUpdateScheduleResponse
	if err := json.Unmarshal(body, &updatedSchedule); err != nil {
		return nil, errors.New("error parsing JSON")
	}

	return &updatedSchedule, nil
}

func (s *ScheduleService) GetScheduleByID(scheduleID string) (*models.TwSchedule, error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("schedule id is required")
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule/"+scheduleID,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()
	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GET /schedule/%s returned status %d: %s", scheduleID, resp.StatusCode, string(body))
	}

	// Parse response
	var schedule models.TwSchedule
	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &schedule, nil
}

func (s *ScheduleService) GetScheduleById(scheduleID string) (*core_dtos.TwScheduleResponse, error) {

	resp, err := dms.CallAPI(
		"GET",
		"/schedule/"+scheduleID,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()
	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("GET /schedule/%s returned status %d: %s", scheduleID, resp.StatusCode, string(body))
	}

	// Parse response
	var schedule core_dtos.TwScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &schedule, nil
}

func (s *ScheduleService) DeleteSchedule(scheduleID string, workspaceUser *models.TwWorkspaceUser) error {
	if scheduleID == "" || scheduleID == "0" || scheduleID == "-1" {
		return fmt.Errorf("schedule not found")
	}
	resp, err := dms.CallAPI(
		"DELETE",
		"/schedule/"+scheduleID+"/workspace_user/"+strconv.Itoa(workspaceUser.ID),
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("can not read response body")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	return nil
}

func (h *ScheduleService) GetSchedulesByBoardColumn(workspaceID string, boardColumnId int) ([]models.TwSchedule, error) {
	if workspaceID == "" || workspaceID == "0" {
		return nil, fmt.Errorf("workspace id is required")
	}
	if boardColumnId == 0 || boardColumnId == -1 {
		return nil, fmt.Errorf("board column id is required")
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule/workspace/"+workspaceID+"/board_column/"+strconv.Itoa(boardColumnId),
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()
	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var schedules []models.TwSchedule
	if err := json.NewDecoder(resp.Body).Decode(&schedules); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return schedules, nil
}

func (s *ScheduleService) GetScheduleDetailByID(scheduleID string) (*models.TwSchedule, error) {

	resp, err := dms.CallAPI(
		"GET",
		"/schedule/"+scheduleID,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()
	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var schedule models.TwSchedule
	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &schedule, nil
}

func (h *ScheduleService) GetSchedulesByBoardColumnWithFilters(workspaceID string, boardColumnId int, filters map[string]interface{}) ([]models.TwSchedule, error) {
	// Construct the base URL
	url := "/schedule/workspace/" + workspaceID + "/board_column/" + strconv.Itoa(boardColumnId) + "/filter"

	// Create queryParams map to hold filter parameters
	queryParams := make(map[string]string)

	// Iterate over the filters map and handle array and boolean types
	for key, value := range filters {
		switch v := value.(type) {
		case string:
			queryParams[key] = v
		case bool:
			// Convert boolean to "true" or "false"
			queryParams[key] = fmt.Sprintf("%v", v)
		case []string:
			// Convert array of strings (e.g., members) into a comma-separated list
			queryParams[key] = strings.Join(v, ",")
		default:
			// For other types, convert to string
			queryParams[key] = fmt.Sprintf("%v", v)
		}
	}

	// Call the API with the constructed URL and query parameters
	resp, err := dms.CallAPI(
		"GET",
		url,
		nil,
		nil,
		queryParams, // Pass queryParams map directly
		120*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("server error: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response into a list of schedules
	var schedules []models.TwSchedule
	if err := json.NewDecoder(resp.Body).Decode(&schedules); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return schedules, nil
}
