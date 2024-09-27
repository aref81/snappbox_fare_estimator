package main

import (
	"fmt"
	"github.com/aref81/snappbox_fare_estimator/hermes/config"
	"github.com/aref81/snappbox_fare_estimator/hermes/internal/delivery"
	"github.com/aref81/snappbox_fare_estimator/hermes/pkg/input/csv"
	"github.com/aref81/snappbox_fare_estimator/shared/broker/rabbitMQ"
	"github.com/aref81/snappbox_fare_estimator/shared/logger"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"go.uber.org/zap"
	"log"
	"os"
	"sync"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	err = logger.InitLogger(zap.InfoLevel)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	zLogger := logger.Logger

	// Initialize RabbitMQ connection
	rabbitMQPublisher, err := rabbitMQ.NewRabbitMQPublisher(cfg.RabbitMQ.URL, cfg.RabbitMQ.Queue, zLogger)
	if err != nil {
		zLogger.Fatal("Failed to initialize RabbitMQ publisher", zap.Error(err))
		return
	}
	defer rabbitMQPublisher.Close()

	deliveryPointChan := make(chan *models.DeliveryPoint, 100)
	wg := sync.WaitGroup{}

	// Initialize reader stream
	reader := csv.NewDeliveryReader(cfg.CSV.FilePath)
	go reader.StreamDeliveryPoints(deliveryPointChan, zLogger)
	wg.Add(1)

	// Initialize publisher stream
	processor := delivery.NewDeliveryProcessor(rabbitMQPublisher, zLogger)
	go processor.ProcessDeliveries(deliveryPointChan)
	wg.Add(1)

	zLogger.Info("Hermes microservice completed successfully")
	wg.Wait()
}
