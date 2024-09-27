package input

import (
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
)

// DeliveryReader is an abstraction for reading Delivery data
type DeliveryReader interface {
	StreamDeliveryPoints(publisherChan chan *models.DeliveryPoint, log *zap.Logger) error
}
