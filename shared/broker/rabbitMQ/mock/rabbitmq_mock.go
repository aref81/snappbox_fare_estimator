package mock

import (
	"errors"
	"github.com/streadway/amqp"
	"sync"
)

// MockRabbitMQ simulates a RabbitMQ broker with multiple Queues
type MockRabbitMQ struct {
	Queues map[string]chan amqp.Delivery
	mutex  sync.Mutex
}

// NewMockRabbitMQ initializes a new mock RabbitMQ broker
func NewMockRabbitMQ() *MockRabbitMQ {
	return &MockRabbitMQ{
		Queues: make(map[string]chan amqp.Delivery),
	}
}

// DeclareQueue declares a new queue with a given name
func (m *MockRabbitMQ) DeclareQueue(queueName string, bufferSize int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.Queues[queueName]; !exists {
		m.Queues[queueName] = make(chan amqp.Delivery, bufferSize)
	}
}

// GetQueue returns the channel for a given queue
func (m *MockRabbitMQ) GetQueue(queueName string) (chan amqp.Delivery, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	queue, exists := m.Queues[queueName]
	if !exists {
		return nil, errors.New("queue not found")
	}
	return queue, nil
}

// CloseQueue closes the specified queue
func (m *MockRabbitMQ) CloseQueue(queueName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	queue, exists := m.Queues[queueName]
	if !exists {
		return errors.New("queue not found")
	}
	close(queue)
	delete(m.Queues, queueName)
	return nil
}
