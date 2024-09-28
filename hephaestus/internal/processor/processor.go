package processor

import (
	"context"
	"encoding/json"
	"github.com/aref81/snappbox_fare_estimator/cmd/hephaestus/pkg/output"
	"github.com/aref81/snappbox_fare_estimator/shared/broker/rabbitMQ"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
	"sync"
)

type Processor struct {
	log                *zap.Logger
	rabbitMQConsumer   *rabbitMQ.RabbitMQConsumer
	deliveryFareWriter output.DeliveryFareWriter
	fareBuffer         []*models.DeliveryFare
	mutex              sync.Mutex
	batchSize          int
}

func NewConsumer(rabbitMQConsumer *rabbitMQ.RabbitMQConsumer, deliveryFareWriter output.DeliveryFareWriter, batchSize int, logger *zap.Logger) (*Processor, error) {
	return &Processor{
		log:                logger,
		rabbitMQConsumer:   rabbitMQConsumer,
		deliveryFareWriter: deliveryFareWriter,
		fareBuffer:         []*models.DeliveryFare{},
		batchSize:          batchSize,
	}, nil
}

func (c *Processor) Consume(ctx context.Context) error {
	msgs, err := c.rabbitMQConsumer.Consume(context.Background())
	if err != nil {
		c.log.Fatal("Failed to consume messages", zap.Error(err))
	}

	for msg := range msgs {
		var fare models.DeliveryFare
		if err := json.Unmarshal(msg.Body, &fare); err != nil {
			c.log.Error("Failed to unmarshal delivery fare", zap.Error(err))
			continue
		}

		c.addToBuffer(&fare)
	}

	<-ctx.Done()
	c.flushBuffer()
	return nil
}

func (c *Processor) addToBuffer(fare *models.DeliveryFare) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.fareBuffer = append(c.fareBuffer, fare)
	if len(c.fareBuffer) >= c.batchSize {
		c.flushBuffer()
	}
}

func (c *Processor) flushBuffer() {
	c.log.Info("Flushing batch to writer", zap.Int("batch_size", len(c.fareBuffer)))

	err := c.deliveryFareWriter.WriteBatch(c.fareBuffer, c.log)
	if err != nil {
		c.log.Error("Failed to write batch", zap.Error(err))
	}

	c.fareBuffer = []*models.DeliveryFare{}
}
