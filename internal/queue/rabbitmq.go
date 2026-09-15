package queue

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial("amqp://taskforge:taskforge@localhost:5672")
	if err != nil {
		return nil, fmt.Errorf("Connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("open RabbitMQ channel: %w", err)
	}

	return &RabbitMQ{
		conn: conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) DeclareQueue() error {
	_, err := r.channel.QueueDeclare(
		"jobs",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("declare jobs queue: %w", err)
	}

	return nil
}