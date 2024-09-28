package main

import (
	"fmt"
	"github.com/aref81/snappbox_fare_estimator/atalanta/config"
	"github.com/aref81/snappbox_fare_estimator/atalanta/internal/processor"
	"github.com/aref81/snappbox_fare_estimator/shared/broker/rabbitMQ"
	"github.com/aref81/snappbox_fare_estimator/shared/logger"
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

	// Initialize RabbitMQ publisher connection
	rabbitMQPublisher, err := rabbitMQ.NewRabbitMQPublisher(cfg.RabbitMQ.URL, cfg.RabbitMQ.FareQueue, zLogger)
	if err != nil {
		zLogger.Fatal("Failed to initialize RabbitMQ publisher", zap.Error(err))
		return
	}
	defer rabbitMQPublisher.Close()

	// Initialize RabbitMQ consumer connection
	rabbitMQConsumer, err := rabbitMQ.NewRabbitMQConsumer(cfg.RabbitMQ.URL, cfg.RabbitMQ.DeliveryQueue, zLogger)
	if err != nil {
		zLogger.Fatal("Failed to initialize RabbitMQ consumer", zap.Error(err))
		return
	}
	defer rabbitMQConsumer.Close()

	wg := sync.WaitGroup{}

	// Initialize prc
	prc := processor.NewProcessor(rabbitMQPublisher, rabbitMQConsumer, zLogger, cfg.FareRules, cfg.TimeBoundaries)
	go prc.ProcessDeliveries()
	wg.Add(1)

	zLogger.Info("Atalanta microservice started successfully")
	wg.Wait()
}
