# Redpanda & mysql cdc Service

A Go-based service that consumes Change Data Capture (CDC) events from MySQL using Debezium connect, Redpanda.

## Overview

This service implements a consumer that listens to database changes (Create, Update, Delete, Read operations) from MySQL and processes them through a Kafka-compatible message broker (Redpanda). It uses Debezium for CDC (Change Data Capture) integration.

## Prerequisites

- Go 1.23.2 or later
- Docker and Docker Compose
- MySQL 8.0
- Redpanda
- Debezium Connect

## Project Structure

```plaintext
.
├── cmd/main.go           # Application entry point
├── consumer/             # Consumer implementation
├── handlers/             # Event handlers
├── scripts/init.sql      # Database schema and seed
├── types/                # Type definitions
├── *-connector.json      # Debezium connector configuration
├── docker-compose.yml    # Docker services configuration
└── go.mod                # Go module file
```

## Setup and Installation

1. Clone the repository

```bash
git clone https://github.com/AbdelilahOu/Redpanda-go.git
```

2. Start the infrastructure services:

```bash
docker-compose up -d
```

3. Create the Debezium connector:

```bash
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @orders-connector.json
```

```bash
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @customers-connector.json
```

4. Run the application:

```bash
go run cmd/main.go
```

## Infrastructure Components

- **MySQL**: Source database (Port: 3306)
- **Redpanda**: Kafka-compatible message broker (Port: 19092)
- **Redpanda Console**: Web UI for managing Redpanda (Port: 8080)
- **Debezium Connect**: CDC connector service (Port: 8083)

## Configuration

### Debezium Connector

The connector is configured to:

- Monitor the `mydb.Orders` table
- Publish events to the `orders.mydb.Orders` topic
- Use JSON format for messages without schemas
- Track schema changes in a dedicated topic

### Consumer Configuration

The Go consumer is configured to:

- Connect to Redpanda broker at `localhost:19092`
- Use consumer group `orders-group`
- Process messages from the `orders.mydb.Orders` topic
- Handle Create, Update, Delete, and Read operations

## Event Handling

The service implements handlers for different database operations:

- Create (`c`): Process new record creation
- Update (`u`): Handle record updates
- Delete (`d`): Process record deletions
- Read (`r`): Handle read events
