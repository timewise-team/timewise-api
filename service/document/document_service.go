package document

import (
	"api/dms"
	"encoding/json"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/document_dtos"
	"github.com/timewise-team/timewise-models/models"
	"strconv"
)

type DocumentService struct {
}

func NewDocumentService() *DocumentService {
	return &DocumentService{}
}
func (h *DocumentService) GetDocumentsBySchedule(scheduleId int) ([]models.TwDocument, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/document/schedule/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var documents []models.TwDocument
	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return nil, err
	}

	return documents, nil
}

func (h *DocumentService) GetDocumentsByScheduleID(scheduleId int) ([]document_dtos.TwDocumentResponse, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/document/schedule_id/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var documents []document_dtos.TwDocumentResponse
	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return nil, err
	}

	return documents, nil
}

func (h *DocumentService) CreateReminder(reminder models.TwReminder) (models.TwReminder, error) {
	reminderJson, err := json.Marshal(reminder)
	if err != nil {
		return models.TwReminder{}, err
	}
	resp, err := dms.CallAPI(
		"POST",
		"/reminder",
		reminderJson,
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

func (h *DocumentService) GetRemindersByScheduleID(id string) (models.TwReminder, error) {
	resp, err := dms.CallAPI("GET", "/reminder/schedule/"+id, nil, nil, nil, 120)
	if err != nil {
		return models.TwReminder{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return models.TwReminder{}, err
	}

	var reminders models.TwReminder
	if err := json.NewDecoder(resp.Body).Decode(&reminders); err != nil {
		return models.TwReminder{}, err
	}

	return reminders, nil
}

func (h *DocumentService) DeleteReminder(id string) error {
	resp, err := dms.CallAPI("DELETE", "/reminder/"+id, nil, nil, nil, 120)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (h *DocumentService) GetReminderByID(id string) interface{} {

}
