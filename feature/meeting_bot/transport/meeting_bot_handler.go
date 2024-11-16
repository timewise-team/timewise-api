package transport

import (
	"api/service/meeting_bot"
	"github.com/gofiber/fiber/v2"
)

type MeetingBotHandler struct {
	service *meeting_bot.MeetingBotService
}

func NewMeetingBotHandler() (*MeetingBotHandler, error) {
	service, err := meeting_bot.NewMeetingBotService()
	if err != nil {
		return nil, err
	}
	return &MeetingBotHandler{
		service: service,
	}, nil
}

type StartMeetingRequest struct {
	ScheduleID int    `json:"schedule_id"`
	MeetLink   string `json:"meet_link"`
}

// StartMeeting godoc
// @Summary Start a meeting
// @Description Start a meeting
// @Security bearerToken
// @Tags meeting_bot
// @Accept json
// @Produce json
// @Param startMeetingRequest body StartMeetingRequest true "Start meeting request"
// @Success 200 {string} string
// @Router /api/v1/meeting_bot/start [post]
func (r MeetingBotHandlerRegister) StartMeeting(ctx *fiber.Ctx) error {
	return r.Handler.service.StartMeeting(ctx)
}
