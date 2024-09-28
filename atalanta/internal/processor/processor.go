package processor

import (
	"context"
	"encoding/json"
	"github.com/aref81/snappbox_fare_estimator/atalanta/config"
	"github.com/aref81/snappbox_fare_estimator/shared/broker/rabbitMQ"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
)

type Processor struct {
	rabbitMQPublisher *rabbitMQ.RabbitMQPublisher
	rabbitMQConsumer  *rabbitMQ.RabbitMQConsumer
	fareCalculator    *fareCalculator
	log               *zap.Logger
}

func NewProcessor(publisher *rabbitMQ.RabbitMQPublisher,
	consumer *rabbitMQ.RabbitMQConsumer,
	log *zap.Logger,
	fareRules config.FareRulesConfig,
	timeBoundaries config.TimeBoundariesConfig) *Processor {
	return &Processor{
		rabbitMQPublisher: publisher,
		rabbitMQConsumer:  consumer,
		fareCalculator: &fareCalculator{
			fareConfig:     fareRules,
			timeBoundaries: timeBoundaries,
		},
		log: log,
	}
}

// ProcessDeliveries consume the deliveries coming from rabbitMQ and process the DeliveryFare for it
func (p *Processor) ProcessDeliveries() {
	msgs, err := p.rabbitMQConsumer.Consume(context.Background())
	if err != nil {
		p.log.Fatal("Failed to consume messages", zap.Error(err))
	}

	for msg := range msgs {
		var delivery models.Delivery
		if err := json.Unmarshal(msg.Body, &delivery); err != nil {
			p.log.Warn("Failed to unmarshal processor", zap.Error(err))
			continue
		}

		go func(delivery *models.Delivery) {
			err := p.processDeliveryFare(delivery)
			if err != nil {
				p.log.Warn("Failed to process Delivery Fare", zap.Error(err))
			}
		}(&delivery)

	}
}

// processDeliveryFare generate the DeliverFare for a single Delivery and push it to the rabbitMQ
func (p *Processor) processDeliveryFare(delivery *models.Delivery) error {
	totalFare := p.fareCalculator.calculateFare(delivery)

	fare := models.DeliveryFare{
		ID:   delivery.ID,
		Fare: totalFare,
	}

	fareBytes, err := json.Marshal(fare)
	if err != nil {
		p.log.Error("Failed to serialize fare", zap.Error(err))
		return err
	}

	err = p.rabbitMQPublisher.PublishMessage(context.Background(), fareBytes)
	if err != nil {
		p.log.Error("Failed to publish fare", zap.Error(err))
		return err
	}

	p.log.Info("Fare calculated and sent", zap.Int("delivery_id", delivery.ID), zap.Float64("total_fare", totalFare))
	return nil
}
