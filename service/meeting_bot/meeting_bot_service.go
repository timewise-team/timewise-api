package meeting_bot

import (
	"api/utils/rabbitmq"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type MeetingBotService struct {
	rabbitMQClient *rabbitmq.RabbitMQClient
}

func NewMeetingBotService() (*MeetingBotService, error) {
	client, err := rabbitmq.NewRabbitMQClient()
	if err != nil {
		return nil, err
	}

	return &MeetingBotService{
		rabbitMQClient: client,
	}, nil
}

func (s *MeetingBotService) StartMeeting(ctx *fiber.Ctx) error {
	var requestBody struct {
		ScheduleID int    `json:"schedule_id"`
		MeetLink   string `json:"meet_link"`
	}

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	message, err := json.Marshal(requestBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to marshal message")
	}

	err = s.rabbitMQClient.Publish("start_meeting_queue", message)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to send message to start_meeting_queue")
	}

	return ctx.SendString("Start meeting")
}
