package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/aref81/snappbox_fare_estimator/hermes/pkg/input"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
	"io"
	"os"
	"strconv"
)

// DeliveryReader implements the I/O interface for working with CSV file
type DeliveryReader struct {
	FilePath string
}

// NewDeliveryReader creates a new CSV reader with the provided file path
func NewDeliveryReader(filePath string) input.DeliveryReader {
	return &DeliveryReader{
		FilePath: filePath,
	}
}

func (r *DeliveryReader) ReadDeliveries(log *zap.Logger) (map[int]*models.Delivery, error) {
	file, err := os.Open(r.FilePath)
	if err != nil {
		log.Error("failed to open file", zap.Error(err))
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	deliveries := make(map[int]*models.Delivery)
	var currentDelivery *models.Delivery

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("Failed to read row", zap.Error(err))
			return nil, fmt.Errorf("failed to read row: %v", err)
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid delivery ID: %v", err)
		}
		lat, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid latitude: %v", err)
		}
		lng, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid longitude: %v", err)
		}
		timestamp, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %v", err)
		}

		if currentDelivery == nil || currentDelivery.ID != id {
			if currentDelivery != nil {
				deliveries[currentDelivery.ID] = currentDelivery
			}
			currentDelivery = models.NewDelivery(id)
		}
		currentDelivery.AddPoint(lat, lng, timestamp)
	}
	if currentDelivery != nil {
		deliveries[currentDelivery.ID] = currentDelivery
	}
	return deliveries, nil
}
