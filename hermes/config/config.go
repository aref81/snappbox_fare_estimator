package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// Config is the structure for the hermes microservice
type Config struct {
	CSVFilePath   string
	RabbitMQURL   string
	RabbitMQQueue string
}

// LoadConfig initializes Viper and loads the configuration from the yaml file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")

	viper.SetDefault("CSVFilePath", "delivery_data.csv")
	viper.SetDefault("RabbitMQURL", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("RabbitMQQueue", "delivery-data")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found, using environment variables and defaults.")
	} else {
		fmt.Printf("Config file loaded: %s\n", viper.ConfigFileUsed())
	}

	logConfig()

	config := &Config{
		CSVFilePath:   viper.GetString("CSVFilePath"),
		RabbitMQURL:   viper.GetString("RabbitMQURL"),
		RabbitMQQueue: viper.GetString("RabbitMQQueue"),
	}

	return config, nil
}

// logConfig prints out all the config values (for debugging)
func logConfig() {
	log.Printf("CSVFilePath: %s", viper.GetString("CSVFilePath"))
	log.Printf("RabbitMQURL: %s", viper.GetString("RabbitMQURL"))
	log.Printf("RabbitMQQueue: %s", viper.GetString("RabbitMQQueue"))
}
