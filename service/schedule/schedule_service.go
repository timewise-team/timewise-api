package schedule

import (
	"api/dms"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
	"github.com/timewise-team/timewise-models/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ScheduleService struct {
}

func NewScheduleService() *ScheduleService {
	return &ScheduleService{}
}

func (s *ScheduleService) CreateSchedule(c *fiber.Ctx, CreateScheduleDto core_dtos.TwCreateScheduleRequest) (interface{}, int, error) {
	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return nil, fiber.StatusInternalServerError, errors.New("Failed to retrieve schedule participant")
	}

	newSchedule := core_dtos.TwCreateScheduleRequest{
		WorkspaceID:     CreateScheduleDto.WorkspaceID,
		BoardColumnID:   CreateScheduleDto.BoardColumnID,
		WorkspaceUserID: &workspaceUser.ID,
		Title:           CreateScheduleDto.Title,
		Description:     CreateScheduleDto.Description,
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

	var result map[string]interface{}
	if err := json.NewDecoder(resp1.Body).Decode(&result); err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return result, resp1.StatusCode, nil
}

func fetchSchedule(scheduleID string) (models.TwSchedule, error) {
	resp, err := dms.CallAPI("GET", "/schedule/ori/"+scheduleID, nil, nil, nil, 120)
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

func applyUpdateFields(baseSchedule, updateSchedule models.TwSchedule, dto core_dtos.TwUpdateScheduleRequest) models.TwSchedule {
	if dto.Title != nil {
		updateSchedule.Title = *dto.Title
	}
	if dto.Description != nil {
		updateSchedule.Description = *dto.Description
	}
	if dto.StartTime != nil {
		// Chuyển chuỗi StartTime thành *time.Time
		parsedStartTime, err := time.Parse("2006-01-02 15:04:05.000", *dto.StartTime)
		if err != nil {
			// Xử lý lỗi nếu không thể phân tích chuỗi thành thời gian
			fmt.Println("Error parsing start time:", err)
		} else {
			updateSchedule.StartTime = &parsedStartTime
		}
	}

	// Chuyển đổi và cập nhật EndTime
	if dto.EndTime != nil {
		// Chuyển chuỗi EndTime thành *time.Time
		parsedEndTime, err := time.Parse("2006-01-02 15:04:05.000", *dto.EndTime)
		if err != nil {
			// Xử lý lỗi nếu không thể phân tích chuỗi thành thời gian
			fmt.Println("Error parsing end time:", err)
		} else {
			updateSchedule.EndTime = &parsedEndTime
		}
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
		updateSchedule.ExtraData = *dto.ExtraData
	}
	if dto.RecurrencePattern != nil {
		updateSchedule.RecurrencePattern = *dto.RecurrencePattern
	}
	if dto.Priority != nil {
		updateSchedule.Priority = *dto.Priority
	}
	return updateSchedule
}

func (s *ScheduleService) UpdateSchedule(c *fiber.Ctx, UpdateScheduleDto core_dtos.TwUpdateScheduleRequest) error {
	scheduleID := c.Params("scheduleId")

	schedule, err := fetchSchedule(scheduleID)
	if err != nil {
		return err
	}

	updateSchedule := schedule
	updateSchedule = applyUpdateFields(schedule, updateSchedule, UpdateScheduleDto)

	scheduleParticipant, ok := c.Locals("scheduleParticipant").(models.TwScheduleParticipant)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve schedule participant",
		})
	}

	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve schedule participant",
		})
	}

	if scheduleParticipant.Status == "assign to" {
		updateSchedule.Title = schedule.Title
		updateSchedule.Description = schedule.Description
		updateSchedule.Location = schedule.Location
		updateSchedule.Visibility = schedule.Visibility
		updateSchedule.VideoTranscript = schedule.VideoTranscript
		updateSchedule.ExtraData = schedule.ExtraData
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
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("cannot read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	var updatedSchedule models.TwSchedule
	if err := json.Unmarshal(body, &updatedSchedule); err != nil {
		return errors.New("error parsing JSON")
	}

	return c.JSON(updatedSchedule)
}

func (s *ScheduleService) UpdateSchedulePosition(scheduleId string, workspaceUser *models.TwWorkspaceUser, UpdateScheduleDto core_dtos.TwUpdateSchedulePosition) (*core_dtos.TwUpdateScheduleResponse, error) {

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

func (s *ScheduleService) GetScheduleById(scheduleID string) (*core_dtos.TwGetScheduleResponse, error) {

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
	var schedule core_dtos.TwGetScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&schedule); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &schedule, nil
}

func (s *ScheduleService) DeleteSchedule(c *fiber.Ctx) error {

	scheduleID := c.Params("scheduleID")
	workspaceUser, ok := c.Locals("workspace_user").(*models.TwWorkspaceUser)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve schedule participant",
		})
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
	// Call API
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
