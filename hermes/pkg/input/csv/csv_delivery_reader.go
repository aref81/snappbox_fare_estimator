package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/aref81/snappbox_fare_estimator/hermes/pkg/input"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
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

// StreamDeliveryPoints reads the CSV file row by row, processes each row, and pushes it to channel
func (r *DeliveryReader) StreamDeliveryPoints(publisherChan chan *models.DeliveryPoint, log *zap.Logger) error {
	file, err := os.Open(r.FilePath)
	if err != nil {
		log.Error("failed to open file", zap.Error(err))
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

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
			log.Warn("Invalid processor ID", zap.String("value", row[0]))
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

		publisherChan <- &models.DeliveryPoint{
			DeliveryID: id,
			Latitude:   lat,
			Longitude:  lng,
			Timestamp:  timestamp,
		}
	}

	log.Info("CSV streaming and publishing completed successfully")
	close(publisherChan)
	return nil
}
