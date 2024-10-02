# Hermes Service

The **Hermes microservice** is responsible for ingesting delivery data (messages) and sending them through the broker, similar to how Hermes delivered messages to the gods.
Hermes reads delivery data from CSV files, processes them into delivery points, and sends these points as messages to the RabbitMQ broker.

---

## Directory Structure

```bash
hermes/
├── Dockerfile
├── README.md
├── cmd
│   └── main.go
├── config
│   └── config.go
├── go.mod
├── go.sum
├── internal
│   └── processor
│       └── processor.go
└── pkg
    └── input
        ├── csv
        │   └── csv_delivery_reader.go
        └── delivery_reader.go
```

---

### **Key Components**

1. **Processor** (`processor.go`)
2. **CSV Delivery Reader** (`csv_delivery_reader.go`)
3. **Configuration** (`config.go`)

---

### **1. Processor (processor.go)**

The `Processor` is responsible for processing incoming delivery points and grouping them into deliveries. Once a delivery is complete, it publishes the delivery information to a message broker (RabbitMQ).

#### Key Elements:
- **Processor Struct**: Contains the `publisher` (which sends data to RabbitMQ) and the logger for logging important information.

- **NewDeliveryProcessor**: Initializes a new `Processor` with a RabbitMQ publisher and a logger.

- **ProcessDeliveries**:
    - This is the core function, which processes incoming delivery points from a channel (`deliveryPointChan`).
    - Each delivery is grouped based on the `DeliveryID`. When all points for a delivery are received, it sends the delivery data to RabbitMQ.
    - If a new delivery starts before the previous one is finished, it processes the previous delivery in the background and starts a new one.
    - The processing of each delivery is handled by `processSingleDelivery`.

- **processSingleDelivery**:
    - Serializes (marshals) the delivery data into JSON and publishes it to RabbitMQ using the publisher.
    - Handles logging for errors or successful processing.

### **2. CSV Delivery Reader (csv_delivery_reader.go)**

The `DeliveryReader` is responsible for reading delivery data from a CSV file and streaming the delivery points to the `Processor`.

#### Key Elements:
- **DeliveryReader Struct**: Holds the file path to the CSV file that contains delivery point data.

- **NewDeliveryReader**: Initializes a `DeliveryReader` by providing the CSV file path.

- **StreamDeliveryPoints**:
    - Opens the CSV file and reads it row by row.
    - Each row is parsed into a `DeliveryPoint`, which contains the delivery ID, latitude, longitude, and timestamp.
    - The delivery points are streamed into a channel (`publisherChan`), which the `Processor` reads from.
    - The function handles error checking for invalid data and logs warnings if any row contains incorrect values (e.g., invalid delivery ID, latitude, longitude, or timestamp).
    - Once all rows are processed, the channel is closed.

### **3. Config (config.go)**

The configuration settings for the Hermes service are defined here. These settings can be loaded from a YAML file or from environment variables.

#### Key Elements:
- **RabbitMQConfig**: Holds RabbitMQ connection details, including the URL and queue name.

- **CSVConfig**: Contains the file path for the CSV file from which the delivery points are read.

- **Config Struct**: Combines the RabbitMQ and CSV configurations.

- **LoadConfig**: Uses the `viper` library to load configuration settings from a YAML file. If no config file is found, it uses environment variables or defaults. It logs the loaded configuration for debugging purposes.

a typical config looks like this:
```yaml
rabbitmq:
  url: "amqp://guest:guest@rabbitmq:5672/"
  queue: "processor-data"

csv:
  file_path: "./data/delivery_data_chunk_0.csv"
```
---

### **Workflow Overview**

1. **Reading Data**: The `DeliveryReader` reads delivery points from a CSV file and pushes them to a channel.
2. **Processing Data**: The `Processor` listens to the channel, groups the points into deliveries, and processes each delivery by publishing it to RabbitMQ.
3. **Publishing to RabbitMQ**: Once a delivery is complete, it is serialized and sent to RabbitMQ using the `processSingleDelivery` function.
4. **Logging**: Both the `DeliveryReader` and `Processor` log important events, such as errors, successful processing, and progress.