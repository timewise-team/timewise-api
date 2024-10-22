package schedule_log

import (
	"api/dms"
	"encoding/json"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_log_dtos"
	"strconv"
)

type ScheduleLogService struct {
}

func NewScheduleLogService() *ScheduleLogService {
	return &ScheduleLogService{}
}

func (h *ScheduleLogService) GetScheduleLogsByScheduleID(scheduleId int) ([]schedule_log_dtos.TwScheduleLogResponse, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule_log/schedule/"+scheduleIdStr,
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
