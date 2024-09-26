package broker

import (
	"context"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
	logger  *zap.Logger
}

// NewRabbitMQPublisher initializes a new RabbitMQ connection and starts a queue
func NewRabbitMQPublisher(brokerURL, queueName string, logger *zap.Logger) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", zap.Error(err))
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to create channel", zap.Error(err))
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error("Failed to declare queue", zap.Error(err))
		return nil, err
	}

	return &RabbitMQPublisher{
		conn:    conn,
		channel: channel,
		queue:   queue,
		logger:  logger,
	}, nil
}

// PublishMessage publishes a new message to the queue
func (p *RabbitMQPublisher) PublishMessage(ctx context.Context, body []byte) error {
	err := p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		p.logger.Error("Failed to publish message to RabbitMQ", zap.Error(err))
		return err
	}

	p.logger.Info("Message published to RabbitMQ", zap.ByteString("message", body))
	return nil
}

// Close closes the RabbitMQ connection.
func (p *RabbitMQPublisher) Close() {
	if err := p.channel.Close(); err != nil {
		p.logger.Error("Failed to close RabbitMQ channel", zap.Error(err))
	}
	if err := p.conn.Close(); err != nil {
		p.logger.Error("Failed to close RabbitMQ connection", zap.Error(err))
	}
}
