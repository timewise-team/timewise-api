package notification

import (
	"api/dms"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"github.com/timewise-team/timewise-models/dtos/core_dtos"
)

func PushNotifications(dto core_dtos.PushNotificationDto) error {
	// Validate required fields
	if dto.UserEmailId == 0 || dto.Type == "" || dto.Message == "" {
		return errors.New("Missing required fields")
	}

	// Get RabbitMQ URL from environment variable
	rabbitMQURL := viper.GetString("RABBITMQ_URL")
	if rabbitMQURL == "" {
		return errors.New("RABBITMQ_URL is not set in the environment variables")
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		errorStr := "Failed to connect to RabbitMQ: " + err.Error()
		return errors.New(errorStr)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		errorStr := "Failed to open a channel: " + err.Error()
		return errors.New(errorStr)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"notification_queue",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		errorStr := "Failed to declare a queue: " + err.Error()
		return errors.New(errorStr)
	}

	// Convert dto to JSON
	body, err := json.Marshal(dto)
	if err != nil {
		errorStr := "Failed to marshal notification data: " + err.Error()
		return errors.New(errorStr)
	}

	// Publish message to the queue
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		errorStr := "Failed to publish a message: " + err.Error()
		return errors.New(errorStr)
	}

	// call dms to insert notification into database
	_, err = dms.CallAPI("POST", "/notification", dto, nil, nil, 120)
	if err != nil {
		return err
	}
	return nil
}
