package input

import (
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
)

// DeliveryReader is an abstraction for reading Delivery data
type DeliveryReader interface {
	ReadDeliveries(log *zap.Logger) (map[int]*models.Delivery, error)
}
