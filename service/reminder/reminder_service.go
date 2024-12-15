package reminder

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
	"time"
)

type ReminderService struct {
}

func NewReminderService() *ReminderService {
	return &ReminderService{}
}

//	func (h *ReminderService) CreateReminderAllParticipant(request CreateReminderRequest) (models.TwReminder, error) {
//		reminder := models.TwReminder{
//			ScheduleId:   request.ScheduleId,
//			ReminderTime: request.ReminderTime,
//		}
//
//		return h.CreateReminder(reminder)
//	}
func (h *ReminderService) CreateReminder(reminder models.TwReminder) (models.TwReminder, error) {

	resp, err := dms.CallAPI(
		"POST",
		"/reminder",
		reminder,
		nil,
		nil,
		120,
	)
	if err != nil {
		return models.TwReminder{}, err
	}
	defer resp.Body.Close()
	var reminderResponse models.TwReminder
	if err := json.NewDecoder(resp.Body).Decode(&reminderResponse); err != nil {
		return models.TwReminder{}, err
	}
	return reminderResponse, nil

}

func (h *ReminderService) GetRemindersByScheduleID(id string) ([]models.TwReminder, error) {
	resp, err := dms.CallAPI("GET", "/reminder/schedule/"+id, nil, nil, nil, 120)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, err
	}

	var reminders []models.TwReminder
	if err := json.NewDecoder(resp.Body).Decode(&reminders); err != nil {
		return nil, err
	}

	return reminders, nil
}

func (h *ReminderService) DeleteReminder(id string) error {
	resp, err := dms.CallAPI("DELETE", "/reminder/"+id, nil, nil, nil, 120)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (h *ReminderService) GetReminderByID(id string) (models.TwReminder, error) {
	resp, err := dms.CallAPI("GET", "/reminder/"+id, nil, nil, nil, 120)
	if err != nil {
		return models.TwReminder{}, err
	}
	defer resp.Body.Close()
	var reminder models.TwReminder
	if err := json.NewDecoder(resp.Body).Decode(&reminder); err != nil {
		return models.TwReminder{}, err
	}
	return reminder, nil

}

func (h *ReminderService) UpdateReminder(id string, request models.TwReminder) error {

	resp, err := dms.CallAPI(
		"PUT",
		"/reminder/"+id,
		request,
		nil,
		nil,
		120,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (h *ReminderService) CreateReminderAllParticipant(scheduleDetail *models.TwSchedule, WorkspaceUser *models.TwWorkspaceUser, ScheduleParticipant models.TwScheduleParticipant, reminderTimeInt int) error {

	var reminderRequests []models.TwReminder
	if WorkspaceUser.Role == "admin" || WorkspaceUser.Role == "owner" || (WorkspaceUser.Role == "member" && ScheduleParticipant.Status == "creator") || (WorkspaceUser.Role == "member" && ScheduleParticipant.Status == "assign to" || (WorkspaceUser.Role == "member" && ScheduleParticipant.Status == "creator")) {
		startTime := scheduleDetail.StartTime
		reminderTime := startTime.Add(-time.Duration(reminderTimeInt) * time.Minute)
		if reminderTime.Before(time.Now()) {
			return errors.New("Reminder time must be in the future")
		}
		//if reminder.Type == "only me" {
		var reminderRequest = models.TwReminder{
			ScheduleId:      scheduleDetail.ID,
			ReminderTime:    reminderTime,
			Type:            "all participants",
			Method:          strconv.Itoa(reminderTimeInt),
			WorkspaceUserID: WorkspaceUser.ID,
			IsSent:          false,
		}
		reminderRequests = append(reminderRequests, reminderRequest)
	} else {
		return errors.New("Forbidden")
	}
	for _, reminderRequest := range reminderRequests {
		_, err := h.CreateReminder(reminderRequest)
		if err != nil {
			return err
		}

	}
	return nil
}
func (h *ReminderService) CreateReminderAllParticipantWhenCreateSchedule(scheduleId int, scheduleStartTime time.Time, WorkspaceUser *models.TwWorkspaceUser, reminderTimeInt int) error {

	var reminderRequests []models.TwReminder
	startTime := scheduleStartTime
	reminderTime := startTime.Add(-time.Duration(reminderTimeInt) * time.Minute)
	//if reminder.Type == "only me" {
	var reminderRequest = models.TwReminder{
		ScheduleId:      scheduleId,
		ReminderTime:    reminderTime,
		Type:            "all participants",
		Method:          strconv.Itoa(reminderTimeInt),
		WorkspaceUserID: WorkspaceUser.ID,
		IsSent:          false,
	}
	reminderRequests = append(reminderRequests, reminderRequest)

	for _, reminderRequest1 := range reminderRequests {
		_, err := h.CreateReminder(reminderRequest1)
		if err != nil {
			return err
		}

	}
	return nil
}
