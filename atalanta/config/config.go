package config

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// RabbitMQConfig holds RabbitMQ connection details
type RabbitMQConfig struct {
	URL           string `mapstructure:"url" json:"url"`
	DeliveryQueue string `mapstructure:"delivery_queue" json:"delivery_queue"`
	FareQueue     string `mapstructure:"fare_queue" json:"fare_queue"`
}

// ServiceConfig holds service configuration details
type ServiceConfig struct {
	Port     int    `mapstructure:"port" json:"port"`
	LogLevel string `mapstructure:"log_level" json:"log_level"`
}

// FareRulesConfig holds fare calculation rules
type FareRulesConfig struct {
	MaxSpeed             float64 `mapstructure:"max_speed" json:"max_speed"`
	MinFare              float64 `mapstructure:"min_fare" json:"min_fare"`
	FlagAmount           float64 `mapstructure:"flag_amount" json:"flag_amount"`
	IdleFarePerHour      float64 `mapstructure:"idle_fare_per_hour" json:"idle_fare_per_hour"`
	MovingDayFarePerKm   float64 `mapstructure:"moving_day_fare_per_km" json:"moving_day_fare_per_km"`
	MovingNightFarePerKm float64 `mapstructure:"moving_night_fare_per_km" json:"moving_night_fare_per_km"`
}

// TimeBoundariesConfig holds time boundaries rules
type TimeBoundariesConfig struct {
	DayStartHour int `mapstructure:"day_start_hour" json:"day_start_hour"`
	DayEndHour   int `mapstructure:"day_end_hour" json:"day_end_hour"`
}

// Config is the config structure of the Atalanta service
type Config struct {
	RabbitMQ       RabbitMQConfig       `mapstructure:"rabbitmq" json:"rabbitmq"`
	Service        ServiceConfig        `mapstructure:"service" json:"service"`
	FareRules      FareRulesConfig      `mapstructure:"fare_rules" json:"fare_rules"`
	TimeBoundaries TimeBoundariesConfig `mapstructure:"time_boundaries" json:"time_boundaries"`
}

// LoadConfig initializes Viper and loads the configuration from the YAML file
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
