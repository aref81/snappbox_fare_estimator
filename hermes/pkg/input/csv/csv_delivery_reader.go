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

// DeliveryReader implements the input interface for working with CSV file
type DeliveryReader struct {
	FilePath string
}

// NewDeliveryReader creates a new CSV reader with the provided file path
func NewDeliveryReader(filePath string) input.DeliveryReader {
	return &DeliveryReader{
		FilePath: filePath,
	}
}

// ReadDeliveriesAtOnce implements the input interface to read Delivery data from a csv file
func (r *DeliveryReader) ReadDeliveriesAtOnce(log *zap.Logger) (map[int]*models.Delivery, error) {
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

// StreamDeliveries reads the CSV file row by row, processes each row, and pushes it to channel
func (r *DeliveryReader) StreamDeliveries(publisherChan chan *models.Delivery, log *zap.Logger) error {
	file, err := os.Open(r.FilePath)
	if err != nil {
		log.Error("failed to open file", zap.Error(err))
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var currentDelivery *models.Delivery

	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Error("Failed to read row", zap.Error(err))
			return err
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			log.Warn("Invalid delivery ID", zap.String("value", row[0]))
			continue
		}
		lat, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			log.Warn("Invalid latitude", zap.String("value", row[1]), zap.Int("delivery_id", id))
			continue
		}
		lng, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			log.Warn("Invalid longitude", zap.String("value", row[2]))
			continue
		}
		timestamp, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			log.Warn("Invalid timestamp", zap.String("value", row[3]))
			continue
		}

		// If delivery ID changes, process the last delivery and start a new one
		if currentDelivery == nil || currentDelivery.ID != id {
			if currentDelivery != nil {
				// publish previous delivery
				publisherChan <- currentDelivery
			}
			// Create new delivery
			currentDelivery = models.NewDelivery(id)
		}
		err = currentDelivery.AddPoint(lat, lng, timestamp)
		if err != nil {
			log.Warn("Failed to add new delivery point", zap.Error(err))
		}
	}

	if currentDelivery != nil {
		publisherChan <- currentDelivery
	}

	log.Info("CSV streaming and publishing completed successfully")
	return nil
}
