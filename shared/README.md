# Shared Modules

The shared folder contains reusable code for different services in the project, encapsulating logic related to message brokering (RabbitMQ), geographical distance calculation (Haversine), logging, and models that represent the delivery system.

## Directory Structure

```bash
shared/
│
├── broker/
│   ├── rabbitMQ/
│   │   ├── mock/
│   │   │   ├── rabbitmq_mock.go
│   │   │   ├── rabbitmqConsumer_mock.go
│   │   │   └── rabbitmqPublisher_mock.go
│   │   ├── rabbitmqConsumer.go
│   │   └── rabbitmqPublisher.go
│   ├── broker.go
│   └── go.mod
├── haversine/
│   ├── haversine.go
│   └── haversine_test.go
├── logger/
│   ├── logger.go
│   └── go.mod
├── models/
│   ├── delivery.go
│   ├── delivery_fare.go
│   └── delivery_test.go
└── README.md
```

---

### **Broker Package**

#### 1. `broker.go`
Defines the interfaces for the message broker:
- **Publisher**: Interface for publishing messages.
- **Consumer**: Generic interface for consuming messages from a queue, using the Go generics (`T any`).

#### 2. `rabbitmqConsumer.go`
Implements a RabbitMQ consumer:
- **RabbitMQConsumer struct**: Handles the connection to RabbitMQ, including the channel and queue setup.
- **NewRabbitMQConsumer function**: Initializes a new RabbitMQ consumer and declares the necessary queue.
- **Consume function**: Consumes messages from the specified RabbitMQ queue.
- **Close function**: Closes the RabbitMQ connection and channel.

#### 3. `rabbitmqPublisher.go`
Implements a RabbitMQ publisher:
- **RabbitMQPublisher struct**: Manages the connection, channel, and queue for publishing messages to RabbitMQ.
- **NewRabbitMQPublisher function**: Initializes a new publisher and declares the queue.
- **PublishMessage function**: Publishes a message to the RabbitMQ queue.
- **Close function**: Closes the RabbitMQ connection and channel.

#### 4. **Mock Directory**
Contains mock implementations for testing the RabbitMQ integration:
- **rabbitmq_mock.go**: A mock implementation for RabbitMQ.
- **rabbitmqConsumer_mock.go**: A mock consumer for testing.
- **rabbitmqPublisher_mock.go**: A mock publisher for testing.

---

### **Haversine Package**

#### 1. `haversine.go`
Implements the **Haversine formula** to calculate the great-circle distance between two points on a sphere, given their latitudes and longitudes:
- **Haversine function**: Takes the latitude and longitude of two points and returns the distance in kilometers.

#### 2. `haversine_test.go`
Unit tests for the Haversine function to ensure correct distance calculations.

---

### **Logger Package**

#### 1. `logger.go`
Configures a global **Zap logger** for structured logging:
- **InitLogger function**: Initializes the logger with a given logging level and output format (JSON), and configures logging to `stdout` and error output to `stderr`.

---

### **Models Package**

#### 1. `delivery.go`
Defines the model for deliveries:
- **DeliveryPoint struct**: Represents a GPS coordinate and timestamp for a delivery.
- **DeliverySegment struct**: Represents a segment of the delivery path, with speed, time, and distance.
- **Delivery struct**: Represents a delivery containing multiple segments.
- **AddSegment function**: Adds a validated segment to the delivery.
- **NewDelivery function**: Initializes a new delivery.

#### 2. `delivery_fare.go`
Defines the model for delivery fare calculation:
- **DeliveryFare struct**: Holds the ID and fare amount calculated for a delivery.
- **NewDeliveryFare function**: Initializes a new delivery fare.

#### 3. `delivery_test.go`
Contains unit tests for the delivery and fare models to ensure validation and calculations are correct.