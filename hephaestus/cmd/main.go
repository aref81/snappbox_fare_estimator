package main

import (
	"context"
	"fmt"
	"github.com/aref81/snappbox_fare_estimator/cmd/hephaestus/config"
	"github.com/aref81/snappbox_fare_estimator/cmd/hephaestus/internal/processor"
	"github.com/aref81/snappbox_fare_estimator/cmd/hephaestus/pkg/output/csv"
	"github.com/aref81/snappbox_fare_estimator/shared/broker/rabbitMQ"
	"github.com/aref81/snappbox_fare_estimator/shared/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Initialize logger
	err = logger.InitLogger(zap.InfoLevel)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	zLogger := logger.Logger

	// Initialize RabbitMQ consumer connection
	rabbitMQConsumer, err := rabbitMQ.NewRabbitMQConsumer(cfg.RabbitMQ.URL, cfg.RabbitMQ.Queue, zLogger)
	if err != nil {
		zLogger.Fatal("Failed to initialize RabbitMQ consumer", zap.Error(err))
		return
	}
	defer rabbitMQConsumer.Close()

	// Create CSV Writer
	csvWriter, err := csv.NewCSVWriter(cfg.CSV.FilePath)
	if err != nil {
		zLogger.Fatal("Failed to initialize CSV writer", zap.Error(err))
	}
	defer csvWriter.Close()

	// Graceful shutdown handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}

	consumer, err := processor.NewConsumer(
		rabbitMQConsumer,
		csvWriter,
		cfg.CSV.BatchSize,
		time.Duration(cfg.CSV.FlushInterval)*time.Second,
		zLogger,
	)

	if err != nil {
		zLogger.Fatal("Failed to Consumer Processor", zap.Error(err))
	}

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			zLogger.Fatal("Error while consuming messages", zap.Error(err))
		}
	}()
	wg.Add(1)

	wg.Wait()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zLogger.Info("Shutting down gracefully")
	cancel()
}
