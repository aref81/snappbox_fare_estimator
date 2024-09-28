package processor

import (
	"context"
	"encoding/json"
	"github.com/aref81/snappbox_fare_estimator/shared/broker/rabbitMQ"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
)

type Processor struct {
	rabbitMQPublisher *rabbitMQ.RabbitMQPublisher
	log               *zap.Logger
}

func NewDeliveryProcessor(rabbitMQPublisher *rabbitMQ.RabbitMQPublisher, log *zap.Logger) *Processor {
	return &Processor{
		rabbitMQPublisher: rabbitMQPublisher,
		log:               log,
	}
}

// ProcessDeliveries process all coming deliveries from a channel
func (p *Processor) ProcessDeliveries(deliveryPointChan <-chan *models.DeliveryPoint) error {
	var currentDelivery *models.Delivery
	var previousPoint *models.DeliveryPoint

	for point := range deliveryPointChan {
		// If processor ID changes, process the last processor and start a new one
		if currentDelivery == nil || currentDelivery.ID != point.DeliveryID {
			if currentDelivery != nil {
				// process previous processor
				go func(delivery *models.Delivery) {
					err := p.processSingleDelivery(delivery)
					if err != nil {
						p.log.Warn("Failed to process processor",
							zap.Int("delivery_id", delivery.ID),
							zap.Error(err))
					}
				}(currentDelivery)
			}
			// Create new processor
			currentDelivery = models.NewDelivery(point.DeliveryID)
			previousPoint = nil
		}
		// in case of first point
		if previousPoint != nil {
			err := currentDelivery.AddSegment(*previousPoint, *point)
			if err != nil {
				// if the new point is invalid, we skip this point and reach to the next
				p.log.Warn("Failed to add new processor segment", zap.Error(err))
				continue
			}
			previousPoint = point
		} else {
			previousPoint = point
		}
	}

	// Process the last Delivery
	go func(delivery *models.Delivery) {
		err := p.processSingleDelivery(currentDelivery)
		if err != nil {
			p.log.Warn("Failed to process processor",
				zap.Int("delivery_id", delivery.ID),
				zap.Error(err))
		}
	}(currentDelivery)

	return nil
}

// processSingleDelivery processes a processor, including validation and pushing
func (p *Processor) processSingleDelivery(delivery *models.Delivery) error {
	deliveryBytes, err := json.Marshal(delivery)
	if err != nil {
		p.log.Error("Failed to serialize processor", zap.Int("delivery_id", delivery.ID), zap.Error(err))
		return err
	}

	err = p.rabbitMQPublisher.PublishMessage(context.Background(), deliveryBytes)
	if err != nil {
		p.log.Error("Failed to publish processor to RabbitMQ", zap.Int("delivery_id", delivery.ID), zap.Error(err))
		return err
	}

	p.log.Info("Delivery processed successfully", zap.Int("delivery_id", delivery.ID))
	return nil
}
