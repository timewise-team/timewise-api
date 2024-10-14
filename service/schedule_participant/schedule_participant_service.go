package schedule_participant

import (
	"api/dms"
	"encoding/json"
	"github.com/timewise-team/timewise-models/dtos/core_dtos/schedule_participant_dtos"
	"strconv"
)

type ScheduleParticipantService struct {
}

func NewScheduleParticipantService() *ScheduleParticipantService {
	return &ScheduleParticipantService{}
}

func (h *ScheduleParticipantService) GetScheduleParticipantsBySchedule(scheduleId int, workspaceId string) ([]schedule_participant_dtos.ScheduleParticipantInfo, error) {
	scheduleIdStr := strconv.Itoa(scheduleId)
	if scheduleIdStr == "" {
		return nil, nil
	}
	if workspaceId == "" {
		return nil, nil
	}
	resp, err := dms.CallAPI(
		"GET",
		"/schedule_participant/workspace/"+workspaceId+"/schedule/"+scheduleIdStr,
		nil,
		nil,
		nil,
		120,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var scheduleParticipants []schedule_participant_dtos.ScheduleParticipantInfo
	if err := json.NewDecoder(resp.Body).Decode(&scheduleParticipants); err != nil {
		return nil, err
	}

	return scheduleParticipants, nil

}
