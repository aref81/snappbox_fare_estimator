# Atalanta Service

This service is responsible for calculating the delivery fare based on specific fare rules and time boundaries. It interacts with RabbitMQ for message passing between microservices.

## Directory Structure

```bash
atalanta/
│
├── cmd/
│   └── main.go
├── config/
│   └── config.go
├── internal/
│   └── processor/
│       ├── calculator.go
│       ├── calculator_test.go
│       └── processor.go
├── Dockerfile
├── go.mod
└── README.md
```

### Main Files

#### 1. **`processor.go`**
- This is the core of the Atalanta service, responsible for consuming delivery messages, calculating fares, and publishing the results.
- **Key Components**:
    - **Processor struct**: Holds the consumer, publisher, fareCalculator, and logger.
    - **ProcessDeliveries function**: Consumes deliveries from RabbitMQ and processes each one concurrently.
    - **processDeliveryFare function**: Calculates the fare for each delivery and publishes the result back to RabbitMQ.

#### 2. **`calculator.go`**
- This file handles the actual fare calculation logic based on the configuration (fare rules and time boundaries).
- **Key Components**:
    - **fareCalculator struct**: Contains configuration details related to fare rules and time boundaries.
    - **calculateFare function**: Implements the fare calculation based on distance, speed, and whether the segment occurs during the day or night.
    - **isDayTime function**: Determines if a given timestamp is during the day or night based on the configuration.

#### 3. **`config.go`**
- This file is responsible for loading the configuration from a YAML file or environment variables.
- **Key Components**:
    - **RabbitMQConfig**: Holds RabbitMQ connection details.
    - **ServiceConfig**: Holds service-level configurations like port and log level.
    - **FareRulesConfig**: Defines rules for fare calculation such as idle fare, minimum fare, and fare per kilometer for day/night.
    - **TimeBoundariesConfig**: Defines the time boundaries for day and night.
    - **LoadConfig function**: Loads configuration using Viper and unmarshals it into the defined structs.

a typical config looks like this:
```yaml
rabbitmq:
  url: "amqp://guest:guest@rabbitmq:5672/"
  delivery_queue: "processor-data"
  fare_queue: "fares-data"

fare_rules:
  min_fare: 3.47
  flag_amount: 1.30
  idle_fare_per_hour: 11.90
  moving_day_fare_per_km: 0.74
  moving_night_fare_per_km: 1.30

time_boundaries:
  day_start_hour: 5
  night_end_hour: 24
```