package output

import (
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
)

// DeliveryFareWriter is an abstraction for efficiently writing DeliveryFare data
type DeliveryFareWriter interface {
	WriteBatch(fares []*models.DeliveryFare, log *zap.Logger) error
	Close()
}
