package transport

import (
	"api/service/notification"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
)

type NotificationHandler struct {
	service notification.NotificationService
}

func NewNotificationHandler() (*NotificationHandler, error) {
	service, err := notification.NewNotificationService()
	if err != nil {
		return nil, err
	}

	return &NotificationHandler{
		service: *service,
	}, nil
}

// PushNotifications godoc
// @Summary Push notifications
// @Description Push notifications
// @Tags notification
// @Accept json
// @Produce json
// @Security bearerToken
// @Param message body core_dtos.PushNotificationDto true "Push notification message"
// @Success 200 {object} core_dtos.PushNotificationDto
// @Router /api/v1/notification/push [post]
func (handler *NotificationHandler) PushNotifications(c *fiber.Ctx) error {
	// Get RabbitMQ channel from the service
	channel, err := handler.service.GetChannel()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to connect to RabbitMQ"})
	}

	// Parse the request body into the PushNotificationDto
	message := new(core_dtos.PushNotificationDto)
	if err := c.BodyParser(message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	// Serialize the message to JSON format
	msgBody, err := json.Marshal(message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to serialize message"})
	}

	// Publish message to the queue with content type as application/json
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        msgBody,
	}

	err = channel.Publish("", notification.QueueName, false, false, msg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send message"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": message.Message, "status": "success"})
}
