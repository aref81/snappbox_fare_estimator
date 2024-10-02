package mock

import (
	"context"
	"github.com/streadway/amqp"
)

// MockRabbitMQPublisher implements the RabbitMQPublisher interface for testing
type MockRabbitMQPublisher struct {
	mockBroker *MockRabbitMQ
	queueName  string
}

// NewMockRabbitMQPublisher initializes a new mock publisher
func NewMockRabbitMQPublisher(mockBroker *MockRabbitMQ, queueName string) *MockRabbitMQPublisher {
	return &MockRabbitMQPublisher{
		mockBroker: mockBroker,
		queueName:  queueName,
	}
}

// PublishMessage publishes a message to the specified queue
func (m *MockRabbitMQPublisher) PublishMessage(ctx context.Context, message []byte) error {
	queue, err := m.mockBroker.GetQueue(m.queueName)
	if err != nil {
		return err
	}
	queue <- amqp.Delivery{
		Body: message,
	}
	return nil
}
