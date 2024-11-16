package rabbitmq

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	HOST     string
	PORT     int
	USERNAME string
	PASSWORD string
}

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQClient() (*RabbitMQClient, error) {
	connStr := viper.GetString("RABBITMQ_URL")
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

func (client *RabbitMQClient) Publish(queueName string, message []byte) error {
	_, err := client.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = client.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	return err
}

func (client *RabbitMQClient) Close() {
	client.channel.Close()
	client.conn.Close()
}
