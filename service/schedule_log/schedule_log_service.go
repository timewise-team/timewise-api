package schedule_log

import (
	"api/dms"
	"encoding/json"
	"errors"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_log_dtos"
)

type ScheduleLogService struct {
}

func NewScheduleLogService() *ScheduleLogService {
	return &ScheduleLogService{}
}

func (h *ScheduleLogService) GetScheduleLogsByScheduleID(scheduleId string) ([]schedule_log_dtos.TwScheduleLogResponse, error) {

	if scheduleId == "" || scheduleId == "0" {
		return nil, errors.New("schedule id is required")
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule_log/schedule/"+scheduleId,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleLogs []schedule_log_dtos.TwScheduleLogResponse
	if err := json.NewDecoder(resp.Body).Decode(&scheduleLogs); err != nil {
		return nil, err
	}

	return scheduleLogs, nil

}
