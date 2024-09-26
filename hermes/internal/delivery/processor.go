package delivery

import (
	"context"
	"encoding/json"
	"github.com/aref81/snappbox_fare_estimator/hermes/pkg/broker"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
)

type Processor struct {
	rabbitMQPublisher *broker.RabbitMQPublisher
	log               *zap.Logger
}

func NewDeliveryProcessor(rabbitMQPublisher *broker.RabbitMQPublisher, log *zap.Logger) *Processor {
	return &Processor{
		rabbitMQPublisher: rabbitMQPublisher,
		log:               log,
	}
}

// ProcessDeliveries process all coming deliveries from a channel
func (p Processor) ProcessDeliveries(deliveryChan <-chan *models.Delivery) error {
	for delivery := range deliveryChan {
		err := p.processSingleDelivery(delivery)
		if err != nil {
			p.log.Warn(err.Error())
		}
	}

	return nil
}

// processSingleDelivery processes a delivery, including validation and pushing
func (p Processor) processSingleDelivery(delivery *models.Delivery) error {
	deliveryBytes, err := json.Marshal(delivery)
	if err != nil {
		p.log.Error("Failed to serialize delivery", zap.Int("delivery_id", delivery.ID), zap.Error(err))
		return err
	}

	err = p.rabbitMQPublisher.PublishMessage(context.Background(), deliveryBytes)
	if err != nil {
		p.log.Error("Failed to publish delivery to RabbitMQ", zap.Int("delivery_id", delivery.ID), zap.Error(err))
		return err
	}

	p.log.Info("Delivery processed successfully", zap.Int("delivery_id", delivery.ID))
	return nil
}
