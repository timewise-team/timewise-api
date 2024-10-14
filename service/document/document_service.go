package document

import (
	"api/dms"
	"encoding/json"
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
