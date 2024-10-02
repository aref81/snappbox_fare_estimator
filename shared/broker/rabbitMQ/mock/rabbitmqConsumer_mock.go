package mock

import (
	"context"
	"github.com/streadway/amqp"
)

// MockRabbitMQConsumer implements the RabbitMQConsumer interface for testing
type MockRabbitMQConsumer struct {
	mockBroker *MockRabbitMQ
	queueName  string
}

// NewMockRabbitMQConsumer initializes a new mock consumer
func NewMockRabbitMQConsumer(mockBroker *MockRabbitMQ, queueName string) *MockRabbitMQConsumer {
	return &MockRabbitMQConsumer{
		mockBroker: mockBroker,
		queueName:  queueName,
	}
}

// Consume consumes messages from the specified queue
func (m *MockRabbitMQConsumer) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	queue, err := m.mockBroker.GetQueue(m.queueName)
	if err != nil {
		return nil, err
	}
	return queue, nil
}
