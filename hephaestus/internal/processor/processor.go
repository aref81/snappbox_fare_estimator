package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aref81/snappbox_fare_estimator/cmd/hephaestus/pkg/output"
	"github.com/aref81/snappbox_fare_estimator/shared/broker"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Processor struct {
	log                *zap.Logger
	consumer           broker.Consumer[amqp.Delivery]
	deliveryFareWriter output.DeliveryFareWriter
	fareBuffer         []*models.DeliveryFare
	mutex              sync.Mutex
	batchSize          int
	flushInterval      time.Duration
}

func NewConsumer(
	consumer broker.Consumer[amqp.Delivery],
	deliveryFareWriter output.DeliveryFareWriter,
	batchSize int, flushInterval time.Duration,
	logger *zap.Logger,
) (*Processor, error) {
	return &Processor{
		consumer:           consumer,
		log:                logger,
		deliveryFareWriter: deliveryFareWriter,
		fareBuffer:         []*models.DeliveryFare{},
		batchSize:          batchSize,
		flushInterval:      flushInterval,
	}, nil
}

// Consume receives the DeliveryFare data from RabbitMQ and writes them into a csv file in efficient way
func (p *Processor) Consume(ctx context.Context) error {
	msgs, err := p.consumer.Consume(context.Background())
	if err != nil {
		p.log.Fatal("Failed to consume messages", zap.Error(err))
	}

	ticker := time.NewTicker(p.flushInterval)
	defer ticker.Stop()

	go func() {
		startTime := time.Now()
		i := 0

		for {
			select {
			case msg := <-msgs:
				go func(delivery amqp.Delivery) {
					var fare models.DeliveryFare
					if err := json.Unmarshal(delivery.Body, &fare); err != nil {
						p.log.Error("Failed to unmarshal delivery fare", zap.Error(err))
						return
					}
					p.addToBuffer(&fare)
				}(msg)

			case <-ticker.C:
				p.log.Info("Flushing buffer due to timeout",
					zap.Int("total processed deliveries", i),
					zap.String("Duration", fmt.Sprintf("%s", time.Now().Sub(startTime))))
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
