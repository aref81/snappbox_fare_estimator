# Hephaestus Service

**Hephaestus** is responsible for consuming delivery fare data from RabbitMQ, buffering the data, and periodically writing it to a CSV file. The service is optimized to handle a large volume of messages by buffering them and only writing to disk in batches, reducing the frequency of disk I/O operations.

---
## Directory Structure

```bash
hephaestus/
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
    └── output
        ├── csv
        │   └── csv_DeliveryDareWriter.go
        └── deliveryFareWriter.go
```

---


### **Key Components**

1. **Processor** (`processor.go`)
2. **CSV Writer** (`csv_DeliveryFareWriter.go`)
3. **Configuration** (`config.go`)

---

### **Processor (`processor.go`)**

The `Processor` struct manages the core functionality of Hephaestus, consuming messages from RabbitMQ, buffering the data, and writing it to disk.

- **Fields:**
    - `log`: For logging operations using the **Zap** logger.
    - `consumer`: A generic consumer that consumes messages from RabbitMQ (of type `amqp.Delivery`).
    - `deliveryFareWriter`: Responsible for writing the buffered delivery fare data into a CSV file.
    - `fareBuffer`: A slice that holds the delivery fare data temporarily.
    - `mutex`: Ensures safe concurrent access to the buffer.
    - `batchSize`: The size of the batch after which the buffer is flushed to disk.
    - `flushInterval`: The time interval after which the buffer is flushed, even if the batch size is not reached.

- **Methods:**
    - `NewConsumer`: Initializes a new `Processor` with the provided consumer, writer, batch size, flush interval, and logger.
    - `Consume`: This method:
        - Continuously receives delivery fare messages from RabbitMQ.
        - Adds them to the buffer.
        - Flushes the buffer either when the batch size is reached or when the flush interval elapses.
        - Stops processing when the context is canceled.
    - `addToBuffer`: Adds a delivery fare to the buffer and triggers a flush if the batch size is exceeded.
    - `flushBuffer`: Writes the contents of the buffer to a CSV file and clears the buffer.

---

### **CSV Writer (`csv_DeliveryFareWriter.go`)**

The CSV writer is responsible for writing batches of delivery fare data into a CSV file.

- **Fields:**
    - `file`: A reference to the opened CSV file.
    - `csvWriter`: The actual CSV writer responsible for writing the rows.
    - `mutex`: Ensures that only one write operation can happen at a time (safeguarding against concurrency issues).

- **Methods:**
    - `NewCSVWriter`: Initializes the CSV writer and opens the specified file. If the file doesn't exist, it will be created.
    - `WriteBatch`: Writes a batch of delivery fare data to the CSV file. Each fare is written as a new row.
    - `Close`: Closes the file when writing is complete.

---

### **Configuration (`config.go`)**

The configuration file loads and manages the service's settings, including RabbitMQ connection details and CSV file settings.

- **Fields:**
    - `RabbitMQConfig`: Holds the RabbitMQ URL and queue name.
    - `CSVConfig`: Contains the CSV file path, batch size, and flush interval.
    - `Config`: The main configuration struct that encapsulates RabbitMQ and CSV configurations.

- **Methods:**
    - `LoadConfig`: Initializes **Viper**, loads the configuration from a YAML file, and unmarshals it into a `Config` struct.
    - `logConfig`: Logs the loaded configuration in a structured JSON format for debugging purposes.

a typical config looks like this:
```yaml
rabbitmq:
  url: "amqp://guest:guest@rabbitmq:5672/"
  queue: "fares-data"

csv:
  file_path: "output/output.csv"
  batch_size: 100
  flush_interval: 30
```
---

### **Flow**

1. **Initialization:**
    - The service is initialized by loading the configuration using `LoadConfig`.
    - A RabbitMQ consumer and a CSV writer are created using the loaded configuration.
    - The `Processor` is initialized with the consumer, writer, batch size, flush interval, and logger.

2. **Message Consumption:**
    - The `Processor` consumes messages from RabbitMQ. Each message contains a `DeliveryFare` object in JSON format.
    - The JSON payload is unmarshaled into a `DeliveryFare` struct and added to the buffer.

3. **Buffer Management:**
    - The buffer temporarily holds the delivery fares. If the number of messages in the buffer reaches the `batchSize` or the `flushInterval` timer expires, the buffer is flushed to disk.

4. **Flushing to CSV:**
    - When the buffer is flushed, all buffered delivery fares are written to the CSV file.
    - The file is kept open for appending, and the buffer is cleared after each flush.

5. **Graceful Shutdown:**
    - If the context is canceled (e.g., service shutdown), the processor ensures that any remaining data in the buffer is flushed to disk before exiting.

---

### **Scalability and Efficiency**
- **Batch Processing**: By processing messages in batches, the system reduces the number of file writes, which improves I/O efficiency, especially for large data volumes.
- **Concurrency**: By using Go routines for message processing and employing a mutex for shared resources, the system supports efficient parallel processing without race conditions.
- **Configuration Flexibility**: The batch size and flush interval can be configured, allowing tuning based on the expected data load and system performance.