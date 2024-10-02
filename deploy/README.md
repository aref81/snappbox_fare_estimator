# SnappBox Fare Estimator Deployment Guide

This guide will walk you through deploying the **SnappBox Fare Estimator** project using **Docker Compose**. The project consists of three main services—**Hermes**, **Atalanta**, and **Hephaestus**—and a **RabbitMQ** message broker. These services are containerized, and Docker Compose is used to orchestrate the deployment.

---

## Directory Structure

```bash
deploy/
├── README.md
├── configs
│   ├── atalanta_config.yaml
│   ├── hephaestus_config.yaml
│   └── hermes_config.yaml
├── data
│   ├── delivery_data.csv
│   └── delivery_data_chunk_0.csv
├── docker-compose.yaml
└── output
    └── output.csv
```

---

## Prerequisites

Before you begin, make sure you have the following installed on your system:
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Step-by-Step Deployment

### 1. Clone the Repository

First, clone the project repository from GitHub and navigate to the `deploy` directory:

```bash
git clone https://github.com/aref81/snappbox_fare_estimator.git
cd snappbox_fare_estimator
cd deploy
```

### 2. Configure Services

Each service (Hermes, Atalanta, Hephaestus) has its configuration files located under the `deploy/configs/` folder. The default configuration files are already set up, but you can customize them if necessary.

- **Hermes**: Configures the CSV data path and RabbitMQ connection details.
- **Atalanta**: Contains the logic for calculating delivery fares.
- **Hephaestus**: Handles storing processed delivery fares in output files.

Make sure to check the README file of each microservice for specific configuration details.

### 3. Add CSV Data

To process delivery data, Hermes expects a CSV file as input. By default, the data should be placed in the `data/` directory.

- The default CSV file path in the Hermes config is set to `data/delivery_data.csv`.
- You can modify this path in the configuration file: `deploy/configs/hermes_config.yaml`.

Ensure that your CSV data is in the correct format and available at the specified location.

### 4. Build and Start the Services

Use **Docker Compose** to build and start the services:

```bash
docker-compose up --build
```

This command will:
- **Start RabbitMQ**: A message broker used to pass messages between services.
- **Build and start Hermes**: Reads the CSV data and publishes it to RabbitMQ.
- **Build and start Atalanta**: Consumes the delivery points from RabbitMQ and calculates the delivery fare.
- **Build and start Hephaestus**: Consumes the calculated fares and writes them to output CSV files.

### *WARNING*: RabbitMQ Startup Delay

RabbitMQ must be fully initialized before the services can start interacting with it. The Docker Compose configuration includes a health check for RabbitMQ and a `10s` wait time to ensure RabbitMQ is up and running before other services start. However, this delay might not always be sufficient, depending on your system. You can adjust the wait time as needed.

### 5. Monitor the Logs

Once the services are up and running, you can monitor their logs to check the progress of data processing:

```bash
docker-compose logs -f
```

Logs will indicate when data processing is complete. For example:
- **Hermes** will log the processing of delivery points from the CSV file.
- **Atalanta** will log the delivery fare calculations.
- **Hephaestus** will log when the output is written to the file.

### 6. Stop the Services

Once all delivery data has been processed and output files have been generated, you can stop the services using the following command:

```bash
docker-compose down
```

This will stop and remove the containers for Hermes, Atalanta, Hephaestus, and RabbitMQ.

---

## Docker Compose Overview

**this section is provided by chatGPT and I overviewed it for potential mistakes**

The `docker-compose.yaml` file is used to orchestrate the deployment. Here's a breakdown of each section:

```yaml
version: '3.8'

services:
  rabbitmq:
    image: "rabbitmq:3-management"
    container_name: rabbitmq
    ports:
      - "5672:5672"        # RabbitMQ messaging port
      - "15672:15672"      # RabbitMQ management web UI
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    healthcheck:
      test: [ "CMD", "rabbitmq-diagnostics", "ping" ]
      interval: 10s
      timeout: 5s
      retries: 5

  hermes:
    build:
      context: ../hermes
      dockerfile: Dockerfile
    container_name: hermes
    volumes:
      - ./data:/root/data  # CSV data file is mounted here
      - ./configs/hermes_config.yaml:/root/config/config.yaml
    depends_on:
      rabbitmq:
        condition: service_healthy

  atalanta:
    build:
      context: ../atalanta
      dockerfile: Dockerfile
    container_name: atalanta
    volumes:
      - ./configs/atalanta_config.yaml:/root/config/config.yaml
    depends_on:
      rabbitmq:
        condition: service_healthy

  hephaestus:
    build:
      context: ../hephaestus
      dockerfile: Dockerfile
    container_name: hephaestus
    volumes:
      - ./configs/hephaestus_config.yaml:/root/config/config.yaml
      - ./output/:/root/output/   # Output directory for processed data
    depends_on:
      rabbitmq:
        condition: service_healthy
```

### Key Services:

1. **RabbitMQ**:
    - The message broker used by the services to communicate.
    - Exposes ports `5672` (for messaging) and `15672` (for RabbitMQ management UI).
    - Includes a health check to ensure RabbitMQ is fully operational before other services start.

2. **Hermes**:
    - Reads the delivery data from the CSV file and pushes delivery points to RabbitMQ.
    - Mounts the local `data/` directory for access to the CSV file.
    - Mounts the Hermes configuration file to set RabbitMQ and CSV paths.

3. **Atalanta**:
    - Consumes delivery points from RabbitMQ and calculates delivery fares.
    - Mounts the Atalanta configuration file.

4. **Hephaestus**:
    - Consumes the calculated fares and writes them to output CSV files.
    - Mounts the output directory for writing processed fare data.
    - Mounts the Hephaestus configuration file.

### Dependencies:

Each service depends on **RabbitMQ** and will only start once RabbitMQ is healthy and ready to accept connections. This ensures that no messages are lost during the startup process.

