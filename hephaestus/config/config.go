package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// RabbitMQConfig holds RabbitMQ connection details
type RabbitMQConfig struct {
	URL   string `mapstructure:"url" json:"url"`
	Queue string `mapstructure:"queue" json:"queue"`
}

// CSVConfig holds CSV file config
type CSVConfig struct {
	FilePath      string `mapstructure:"file_path" json:"file_path"`
	BatchSize     int    `mapstructure:"batch_size" json:"batch_size"`
	FlushInterval int    `mapstructure:"flush_interval" json:"flush_interval"`
}

// Config is the config structure of the Hephaestus service
type Config struct {
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq" json:"rabbitmq"`
	CSV      CSVConfig      `mapstructure:"csv" json:"csv"`
}

// LoadConfig initializes Viper and loads the configuration from the yaml
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found, using environment variables and defaults.")
	} else {
		fmt.Printf("Config file loaded: %s\n", viper.ConfigFileUsed())
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode into config struct: %v", err)
	}

	logConfig(config)

	return config, nil
}

// logConfig prints out all the config values (for debugging)
func logConfig(config *Config) {
	conf, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return
	}
	log.Printf("Hermes Config: %s", conf)
}
