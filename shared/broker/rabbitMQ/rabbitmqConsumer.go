package rabbitMQ

import (
	"context"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	log        *zap.Logger
}

// NewRabbitMQConsumer initializes a new RabbitMQ consumer connection and starts a queue
func NewRabbitMQConsumer(url, queueName string, log *zap.Logger) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ", zap.Error(err))
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open a channel", zap.Error(err))
		return nil, err
	}

	// Declare the queue in case it doesn't already exist
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("Failed to declare the queue", zap.Error(err))
		return nil, err
	}

	return &RabbitMQConsumer{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
		log:        log,
	}, nil
}

// Consume listens to messages from the queue
func (c *RabbitMQConsumer) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.log.Error("Failed to register a consumer", zap.Error(err))
		return nil, err
	}

	return msgs, nil
}

// Close cleans up the connection and channel
func (c *RabbitMQConsumer) Close() {
	if err := c.channel.Close(); err != nil {
		c.log.Error("Failed to close RabbitMQ channel", zap.Error(err))
	}
	if err := c.connection.Close(); err != nil {
		c.log.Error("Failed to close RabbitMQ connection", zap.Error(err))
	}
}
