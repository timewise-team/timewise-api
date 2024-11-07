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

	newSchedule := models.TwSchedule{
		WorkspaceId:   *CreateScheduleDto.WorkspaceID,
		BoardColumnId: *CreateScheduleDto.BoardColumnID,
		Title:         *CreateScheduleDto.Title,
		CreatedBy:     workspaceUser.ID,
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

func applyUpdateFields(baseSchedule, updateSchedule models.TwSchedule, dto core_dtos.TwUpdateScheduleRequest) models.TwSchedule {
	//if dto.BoardColumnID != nil {
	//	updateSchedule.BoardColumnId = *dto.BoardColumnID
	//}
	if dto.Title != nil {
		updateSchedule.Title = *dto.Title
	}
	if dto.Description != nil {
		updateSchedule.Description = *dto.Description
	}
	if dto.StartTime != nil {
		t := *dto.StartTime
		updateSchedule.StartTime = &t
	}
	if dto.EndTime != nil {
		t := *dto.EndTime
		updateSchedule.EndTime = &t
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
	//if dto.Position != nil {
	//	updateSchedule.Position = *dto.Position
	//}
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
		updateSchedule.Position = schedule.Position
		updateSchedule.Priority = schedule.Priority

		//if UpdateScheduleDto.BoardColumnID != nil {
		//	updateSchedule.BoardColumnId = *UpdateScheduleDto.BoardColumnID
		//}
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

func (s *ScheduleService) GetScheduleByID(c *fiber.Ctx) (*models.TwSchedule, error) {

	scheduleID := c.Params("scheduleId")
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
