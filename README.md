# Fare Estimator
This is a Project implemented as a challenge defined by [SnappBox](https://snapp-box.com), the main aim of the project is to calculate the fare of deliveries, based on the delivery points data provided in a CSV file.  
The main challenge is to process gigabytes of data efficiently. To address this, I designed a microservice architecture with 3 independent microservices, each focused on a specific part of the processing pipeline. These services communicate with each other via a broker ([RabbitMQ](https://www.rabbitmq.com)). The architecture is explained in detail later.
The main doc of challenge is available in [SnappBox_Challenge](./docs/).

<img height="200" src="./docs/snappbox-logo.svg" width="200"/>

## Services
The architecture is composed of 3 microservices and a broker:

### Hermes
**Hermes** is the Greek god of messengers, travelers, and communication. The **Hermes microservice** is responsible for ingesting delivery data (messages) and sending them through the broker, similar to how Hermes delivered messages to the gods.  
Hermes reads delivery data from CSV files, processes them into delivery points, and sends these points as messages to the RabbitMQ broker.

[More details on Hermes can be found here](./hermes/README.md).

### Atalanta
**Atalanta** is a Greek heroine known for her incredible speed and prowess in running races. This service processes data quickly and efficiently, calculating the fare for each delivery, akin to Atalanta's legendary speed.  
Atalanta consumes delivery data from the broker, processes the data by calculating delivery fares based on time and distance, and then sends the calculated fares to RabbitMQ for further processing.

[More details on Atalanta can be found here](./atalanta/README.md).

### Hephaestus
**Hephaestus** is the Greek god of craftsmanship, metalworking, and creation. The **Hephaestus service** crafts and outputs the final fare data, generating reports in the form of CSV files, similar to how Hephaestus forged intricate creations.  
Hephaestus reads the calculated fares from RabbitMQ and writes them to CSV files in batches, ensuring efficient processing even for large datasets.

[More details on Hephaestus can be found here](./hephaestus/README.md).

### Broker (RabbitMQ)
RabbitMQ is the broker that facilitates communication between the services. Each service operates independently but exchanges messages using queues, allowing for scalable and decoupled communication.

---

## Shared Modules
The shared modules contain the common logic and utility functions used across the services. These modules include:

### broker
This module provides an abstraction layer for RabbitMQ producers and consumers. It allows services to publish and consume messages from RabbitMQ queues without directly handling RabbitMQ configurations, making the interaction with the broker more straightforward and consistent across services.

### haversine
The **Haversine** module calculates the distance between two geographical points using their latitude and longitude. This function is crucial for determining the distance between two delivery points and, in turn, calculating delivery fares.

### logger
The **logger** module provides structured logging using the **Zap** library from Uber. It is configured to provide consistent and structured logging across all services, aiding in debugging and monitoring.

### models
The **models** module contains the data structures used across the services. These include the `Delivery`, `DeliveryPoint`, and `DeliveryFare` models, which represent the core data passed between services during the fare estimation process.


[More details on Shared Modules can be found here](./shared/README.md)

---

## Microservice Structure
Each microservice is designed to focus on a single responsibility within the fare estimation process:

1. **Hermes**: Reads the delivery data from CSV files as a stream and sends delivery points to RabbitMQ for **Atalanta**.
2. **Atalanta**: Consumes delivery points, calculates fares, and sends the calculated fares to RabbitMQ for **Hephaestus**.
3. **Hephaestus**: Consumes calculated fares from RabbitMQ and writes the output to a CSV file by batching them for efficiency.

This structure provides separation of concerns, making it easier to scale and maintain each part of the system independently.

---

## Challenges Faced

### Handling Large Volumes of Data
The biggest challenge was processing gigabytes of delivery data efficiently.
The solution was to design a microservice architecture where each service handles a specific task.
By utilizing RabbitMQ, the data could be streamed between services,
allowing each one to process data independently and avoid bottlenecks.
This structure helped to process data asynchronous, and by use of *goroutines(light abstraction of threads in golang)*
services could process each unit of data concurrently and use the best of resources.

### Ensuring Scalability
To ensure scalability, each microservice can be scaled horizontally as needed.
RabbitMQ acts as a decoupling layer, allowing services to consume and produce messages at their own pace,
which prevents one service from overwhelming another. by using this and keeping the services **stateless**, 
this architecture may be scaled as intended to perform better under great loads of data.

### Efficient File Writing
Since Hephaestus writes large amounts of data to CSV files,
batch writing was implemented to minimize file I/O operations.
The system writes data in batches rather than writing each fare individually, improving performance significantly.
the parameters of batch writing can be configured to fit the pace of data. there is also a timer flush
to avoid keeping data for too long in the buffer.

---

## Deployment Instructions

The project can be deployed using **Docker** and **Docker Compose**. The services are containerized and can be orchestrated easily using the following steps:

1. **Clone the repository**:
    ```bash
   git clone https://github.com/aref81/snappbox_fare_estimator.git
   cd snappbox_fare_estimator
   cd depoly 
   ```
   
2. **Configure services**
    Check the configuration for each service in the mounted locations in `docker-compose.yaml`
    the default config file for each service is under the `deploy/configs` path. a description of configuration
    for each service is available in its own `README.md`, which is mentioned below.

3. **Add data**:
   Place your CSV delivery data in the mounted volume specified in docker compose.
   Make sure to check the config of **Hermes** for loading data correctly. The default path is
   `data/delivery_data.csv` and `data/` is mounted.
   Hermes will read this file, process the data, and begin the pipeline.

4. **Build and start the services** using Docker Compose:
    ```bash
   docker-compose up --build
    ```

This command will:
- Start RabbitMQ.
- Build and start all three microservices (Hermes, Atalanta, and Hephaestus).

Also note that rabbitMQ must be ready before the services start, the config indicates a healthcheck
and a `10s` wait time to let the rabbitMQ start functioning, but as docker health check may not work
properly, config the wait time so that all services are run after the environment is ready.

    
5. **Stop services**:
   After Processing is finished (which is indicated by the logs) to stop the services, run:
    ```bash
   docker-compose down
    ```

3. **Add data**:
   Place your CSV delivery data in the `data/delivery_data.csv` file. Hermes will read this file, process the data, and begin the pipeline.

4. **Stop services**:
   To stop the services, run:


[For a better Guideline on running the project see this](./deploy/README.md)

---
At last the service (as the logs suggest) was able to process about `1000000 Deliveries`, each containing about 
`20 DeliveryPoints` in about `59 seconds`.  

---

## Testing the Services

The project includes unit and end-to-end tests to verify the functionality of each service. These test files ensure that the core logic of each microservice works as expected, both independently and as part of the larger system. Below is an explanation of how the tests are structured and what they cover.

### Test Files Overview

- **Hermes**:
    - The test files for **Hermes** validate the process of reading the CSV file, parsing delivery points, and sending the correct messages to RabbitMQ.
    - These tests focus on ensuring that the correct data is extracted and passed into the pipeline as expected.

- **Atalanta**:
    - The tests for **Atalanta** validate the logic for calculating delivery fares based on the input delivery points consumed from RabbitMQ.
    - They verify that the calculation algorithms work correctly under various conditions, including edge cases like missing or invalid data.

- **Hephaestus**:
    - The test files for **Hephaestus** validate the process of consuming calculated fares and writing them to the output CSV files.
    - These tests ensure that the system can handle large volumes of data and writes them correctly to the expected output.

### Mocking RabbitMQ for Tests

To isolate the testing of business logic from the actual message-broker communication, **RabbitMQ** is mocked during tests. This allows each service to be tested independently without needing to start an actual RabbitMQ instance.

- **Hermes**: When testing Hermes, instead of publishing real messages to RabbitMQ, the RabbitMQ interactions are mocked. This ensures that the message-passing logic works as expected without relying on an external service.

- **Atalanta**: Similarly, for Atalanta, incoming delivery points are simulated via a RabbitMQ mock, which allows the fare calculation logic to be tested in isolation.

- **Hephaestus**: The mock is also used to simulate receiving calculated fares, so the file-writing logic can be tested without needing to interact with a live RabbitMQ instance.

By using this mock setup, the tests can run faster and more reliably in any environment, without the need for external dependencies like RabbitMQ.

### Running the Tests

To run the tests, execute the following command with the appropriate path:

```bash
cd <service_directory>
go test ./...
```

This will run all the unit and end-to-end tests across the services. Make sure you have the required testing dependencies installed before running the tests.

---

## Directory Tree
The overall project structure is as follows:

```bash
snappbox_fare_estimator/
├── LICENSE
├── README.md
├── atalanta
│   ├── Dockerfile
│   ├── README.md
│   ├── cmd
│   │   └── main.go
│   ├── config
│   │   └── config.go
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       └── processor
│           ├── calculator.go
│           ├── calculator_test.go
│           └── processor.go
├── deploy
│   ├── README.md
│   ├── configs
│   │   ├── atalanta_config.yaml
│   │   ├── hephaestus_config.yaml
│   │   └── hermes_config.yaml
│   ├── data
│   │   ├── delivery_data.csv
│   │   └── delivery_data_chunk_0.csv
│   ├── docker-compose.yaml
│   └── output
│       └── output.csv
├── hephaestus
│   ├── Dockerfile
│   ├── README.md
│   ├── cmd
│   │   └── main.go
│   ├── config
│   │   └── config.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   └── processor
│   │       └── processor.go
│   └── pkg
│       └── output
│           ├── csv
│           │   └── csv_DeliveryDareWriter.go
│           └── deliveryFareWriter.go
├── hermes
│   ├── Dockerfile
│   ├── README.md
│   ├── cmd
│   │   └── main.go
│   ├── config
│   │   └── config.go
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   └── processor
│   │       └── processor.go
│   └── pkg
│       └── input
│           ├── csv
│           │   └── csv_delivery_reader.go
│           └── delivery_reader.go
└── shared
    ├── README.md
    ├── broker
    │   ├── broker.go
    │   ├── go.mod
    │   ├── go.sum
    │   └── rabbitMQ
    │       ├── mock
    │       │   ├── rabbitmqConsumer_mock.go
    │       │   ├── rabbitmqPublisher_mock.go
    │       │   └── rabbitmq_mock.go
    │       ├── rabbitmqConsumer.go
    │       └── rabbitmqPublisher.go
    ├── haversine
    │   ├── go.mod
    │   ├── go.sum
    │   ├── haversine.go
    │   └── haversine_test.go
    ├── logger
    │   ├── go.mod
    │   ├── go.sum
    │   └── logger.go
    └── models
        ├── delivery.go
        ├── delivery_fare.go
        ├── delivey_test.go
        ├── go.mod
        └── go.sum
```

---
For further details on each service, please refer to their individual README files:

- [Hermes](./hermes/README.md)
- [Atalanta](./atalanta/README.md)
- [Hephaestus](./hephaestus/README.md)