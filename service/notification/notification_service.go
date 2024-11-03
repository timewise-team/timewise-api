package notification

import (
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

const QueueName = "notifications"

type NotificationService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewNotificationService initializes a persistent RabbitMQ connection and channel
func NewNotificationService() (*NotificationService, error) {
	// Get RabbitMQ URL from environment variable
	rabbitMQURL := os.Getenv("RABBITMQ_URL")

	// Establish RabbitMQ connection
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	// Open channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare the queue
	_, err = channel.QueueDeclare(QueueName, false, false, false, false, nil)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &NotificationService{conn: conn, channel: channel}, nil
}

// GetChannel returns the existing RabbitMQ channel
func (service *NotificationService) GetChannel() (*amqp.Channel, error) {
	if service.channel == nil {
		return nil, errors.New("RabbitMQ channel is not initialized")
	}
	return service.channel, nil
}

// Close closes the RabbitMQ connection and channel
func (service *NotificationService) Close() {
	if service.channel != nil {
		service.channel.Close()
	}
	if service.conn != nil {
		service.conn.Close()
	}
}
