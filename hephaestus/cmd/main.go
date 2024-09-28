package main

import (
	"context"
	"encoding/csv"
	"github.com/aref81/snappbox_fare_estimator/cmd/hephaestus/internal/processor"
	csv2 "github.com/aref81/snappbox_fare_estimator/cmd/hephaestus/pkg/output/csv"
	"github.com/aref81/snappbox_fare_estimator/hephaestus/config"
	"github.com/aref81/snappbox_fare_estimator/shared/broker/rabbitMQ"
	"github.com/aref81/snappbox_fare_estimator/shared/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	// Initialize Logger
	zLogger, err := logger.NewLogger(cfg.Logger.Level)
	if err != nil {
		panic(err)
	}
	defer zLogger.Sync()

	// Initialize RabbitMQ consumer connection
	rabbitMQConsumer, err := rabbitMQ.NewRabbitMQConsumer(cfg.RabbitMQ.URL, cfg.RabbitMQ.DeliveryQueue, zLogger)
	if err != nil {
		zLogger.Fatal("Failed to initialize RabbitMQ consumer", zap.Error(err))
		return
	}
	defer rabbitMQConsumer.Close()

	// Create CSV Writer
	csvWriter, err := csv2.NewCSVWriter(cfg.CSV.FilePath)
	if err != nil {
		zLogger.Fatal("Failed to initialize CSV writer", zap.Error(err))
	}
	defer csvWriter.Close()

	// Graceful shutdown handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}

	consumer, err := processor.NewConsumer(rabbitMQConsumer, csvWriter, cfg.CSV.BatchSize, zLogger)
	if err != nil {
		zLogger.Fatal("Failed to Consumer Processor", zap.Error(err))
	}

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			zLogger.Fatal("Error while consuming messages", zap.Error(err))
		}
	}()
	wg.Add(1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zLogger.Info("Shutting down gracefully")
	wg.Wait()
	cancel()
}
