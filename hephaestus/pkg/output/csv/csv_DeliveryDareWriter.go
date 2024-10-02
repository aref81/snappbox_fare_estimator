package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/aref81/snappbox_fare_estimator/cmd/hephaestus/pkg/output"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"os"
	"sync"

	"go.uber.org/zap"
)

type Writer struct {
	file      *os.File
	csvWriter *csv.Writer
	mutex     sync.Mutex
}

func NewCSVWriter(filePath string) (output.DeliveryFareWriter, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	csvWriter := csv.NewWriter(file)

	return &Writer{
		file:      file,
		csvWriter: csvWriter,
	}, nil
}

func (w *Writer) WriteBatch(fares []*models.DeliveryFare, log *zap.Logger) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, fare := range fares {
		row := []string{
			fmt.Sprintf("%d", fare.ID),
			fmt.Sprintf("%.2f", fare.Fare),
		}
		if err := w.csvWriter.Write(row); err != nil {
			log.Error("Failed to write to CSV", zap.Error(err))
			return err
		}
	}

	w.csvWriter.Flush()
	log.Info("Batch successfully written to CSV")
	return nil
}

func (w *Writer) Close() {
	w.file.Close()
}
