package input

import (
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
)

// DeliveryReader is an abstraction for reading Delivery data
type DeliveryReader interface {
	ReadDeliveriesAtOnce(log *zap.Logger) (map[int]*models.Delivery, error)
	StreamDeliveries(publisherChan chan *models.Delivery, log *zap.Logger) error
}
