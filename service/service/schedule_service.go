package service

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
	"strings"
)

type ScheduleService struct {
}

func NewScheduleService() *ScheduleService {
	return &ScheduleService{}
}

func (s *ScheduleService) CreateSchedule(c *fiber.Ctx, CreateScheduleDto core_dtos.TwCreateScheduleRequest) error {

	workspaceUserIdStr := strconv.Itoa(*CreateScheduleDto.WorkspaceUserID)

	resp, err := dms.CallAPI(
		"GET",
		"/workspace_user/"+workspaceUserIdStr,
		nil,
		nil,
		nil,
		120,
	)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return c.Status(resp.StatusCode).SendString(string(body))
	}

	var workspaceUser models.TwWorkspaceUser
	if err := json.NewDecoder(resp.Body).Decode(&workspaceUser); err != nil {
		return err
	}

	allowedRoles := map[string]bool{
		"admin":  true,
		"owner":  true,
		"member": true,
	}

	if !allowedRoles[strings.ToLower(workspaceUser.Role)] {
		return errors.New("You have no permission to create schedule")
	}

	newSchedule := models.TwSchedule{
		WorkspaceId:       *CreateScheduleDto.WorkspaceID,
		BoardColumnId:     *CreateScheduleDto.BoardColumnID,
		Title:             *CreateScheduleDto.Title,
		Description:       *CreateScheduleDto.Description,
		StartTime:         *CreateScheduleDto.StartTime,
		EndTime:           *CreateScheduleDto.EndTime,
		Location:          *CreateScheduleDto.Location,
		CreatedBy:         *CreateScheduleDto.WorkspaceUserID,
		AllDay:            *CreateScheduleDto.AllDay,
		Status:            *CreateScheduleDto.Status,
		RecurrencePattern: *CreateScheduleDto.RecurrencePattern,
		Visibility:        *CreateScheduleDto.Visibility,
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
		return err
	}
	defer resp1.Body.Close()

	body, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return err
	}

	return c.Status(resp1.StatusCode).Send(body)
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

func fetchWorkspaceUser(workspaceUserIdStr string) (models.TwWorkspaceUser, error) {
	resp, err := dms.CallAPI("GET", "/workspace_user/"+workspaceUserIdStr, nil, nil, nil, 120)
	if err != nil {
		return models.TwWorkspaceUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return models.TwWorkspaceUser{}, fmt.Errorf("GET /workspace_user/%s returned status %d: %s", workspaceUserIdStr, resp.StatusCode, string(body))
	}

	var workspaceUser models.TwWorkspaceUser
	if err := json.NewDecoder(resp.Body).Decode(&workspaceUser); err != nil {
		return models.TwWorkspaceUser{}, fmt.Errorf("error decoding workspace_user response: %v", err)
	}

	return workspaceUser, nil
}

func fetchScheduleParticipant(workspaceUserIdStr, scheduleID string) (models.TwScheduleParticipant, error) {
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
	if dto.BoardColumnID != nil {
		updateSchedule.BoardColumnId = *dto.BoardColumnID
	}
	if dto.Title != nil {
		updateSchedule.Title = *dto.Title
	}
	if dto.Description != nil {
		updateSchedule.Description = *dto.Description
	}
	if dto.StartTime != nil {
		updateSchedule.StartTime = *dto.StartTime
	}
	if dto.EndTime != nil {
		updateSchedule.EndTime = *dto.EndTime
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
	return updateSchedule
}

func (s *ScheduleService) UpdateSchedule(c *fiber.Ctx, UpdateScheduleDto core_dtos.TwUpdateScheduleRequest) error {
	scheduleID := c.Params("scheduleId")
	workspaceUserIdStr := strconv.Itoa(*UpdateScheduleDto.WorkspaceUserID)

	schedule, err := fetchSchedule(scheduleID)
	if err != nil {
		return err
	}

	workspaceUser, err := fetchWorkspaceUser(workspaceUserIdStr)
	if err != nil {
		return err
	}

	updateSchedule := schedule

	if schedule.WorkspaceId == workspaceUser.WorkspaceId && (strings.ToLower(workspaceUser.Role) == "admin" || strings.ToLower(workspaceUser.Role) == "owner") {
		updateSchedule = applyUpdateFields(schedule, updateSchedule, UpdateScheduleDto)
	} else {
		scheduleParticipant, err := fetchScheduleParticipant(workspaceUserIdStr, scheduleID)
		if err != nil {
			return err
		}

		switch strings.ToLower(scheduleParticipant.Status) {
		case "creator":
			updateSchedule = applyUpdateFields(schedule, updateSchedule, UpdateScheduleDto)
		case "assign to":
			updateSchedule.Title = schedule.Title
			updateSchedule.Description = schedule.Description
			updateSchedule.Location = schedule.Location
			updateSchedule.Visibility = schedule.Visibility
			updateSchedule.VideoTranscript = schedule.VideoTranscript
			updateSchedule.ExtraData = schedule.ExtraData
			updateSchedule.RecurrencePattern = schedule.RecurrencePattern
			updateSchedule.StartTime = schedule.StartTime
			updateSchedule.EndTime = schedule.EndTime

			if UpdateScheduleDto.BoardColumnID != nil {
				updateSchedule.BoardColumnId = *UpdateScheduleDto.BoardColumnID
			}
			if UpdateScheduleDto.Status != nil {
				updateSchedule.Status = *UpdateScheduleDto.Status
			}
		default:
			return c.Status(fiber.StatusForbidden).SendString("Forbidden")
		}
	}

	resp, err := dms.CallAPI("PUT", "/schedule/"+scheduleID, updateSchedule, nil, nil, 120)
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

func (s *ScheduleService) DeleteSchedule(c *fiber.Ctx) error {

	scheduleID := c.Params("scheduleID")

	resp, err := dms.CallAPI(
		"DELETE",
		"/schedule/"+scheduleID,
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
