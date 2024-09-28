package processor

import (
	"context"
	"encoding/json"
	"github.com/aref81/snappbox_fare_estimator/cmd/hephaestus/pkg/output"
	"github.com/aref81/snappbox_fare_estimator/shared/broker/rabbitMQ"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Processor struct {
	log                *zap.Logger
	rabbitMQConsumer   *rabbitMQ.RabbitMQConsumer
	deliveryFareWriter output.DeliveryFareWriter
	fareBuffer         []*models.DeliveryFare
	mutex              sync.Mutex
	batchSize          int
	flushInterval      time.Duration
}

func NewConsumer(
	rabbitMQConsumer *rabbitMQ.RabbitMQConsumer,
	deliveryFareWriter output.DeliveryFareWriter,
	batchSize int, flushInterval time.Duration,
	logger *zap.Logger,
) (*Processor, error) {
	return &Processor{
		log:                logger,
		rabbitMQConsumer:   rabbitMQConsumer,
		deliveryFareWriter: deliveryFareWriter,
		fareBuffer:         []*models.DeliveryFare{},
		batchSize:          batchSize,
		flushInterval:      flushInterval,
	}, nil
}

// Consume receives the DeliveryFare data from RabbitMQ and writes them into a csv file in efficient way
func (p *Processor) Consume(ctx context.Context) error {
	msgs, err := p.rabbitMQConsumer.Consume(context.Background())
	if err != nil {
		p.log.Fatal("Failed to consume messages", zap.Error(err))
	}

	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case msg := <-msgs:
				var fare models.DeliveryFare
				if err := json.Unmarshal(msg.Body, &fare); err != nil {
					p.log.Error("Failed to unmarshal delivery fare", zap.Error(err))
					continue
				}
				p.addToBuffer(&fare)

			case <-ticker.C:
				p.log.Info("Flushing buffer due to timeout")
				p.flushBuffer()

			case <-ctx.Done():
				p.flushBuffer()
				return
			}
		}
	}()

	<-ctx.Done()
	return nil
}

// addToBuffer temporarily stores the received Deliveries in memory to later flush them into disk
// flushing is done when batch size exceeds or timeout is reached
func (p *Processor) addToBuffer(fare *models.DeliveryFare) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.fareBuffer = append(p.fareBuffer, fare)
	if len(p.fareBuffer) >= p.batchSize {
		p.flushBuffer()
	}
}

// flushBuffer writes the data in the buffer into a csv file on disk
func (p *Processor) flushBuffer() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.fareBuffer) == 0 {
		return
	}

	p.log.Info("Flushing batch to writer", zap.Int("batch_size", len(p.fareBuffer)))

	err := p.deliveryFareWriter.WriteBatch(p.fareBuffer, p.log)
	if err != nil {
		p.log.Error("Failed to write batch", zap.Error(err))
	}

	p.fareBuffer = []*models.DeliveryFare{}
}
