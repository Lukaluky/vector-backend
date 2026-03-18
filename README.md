# Shipment Service

gRPC-based backend service for managing shipments and tracking their status history. Built with Clean Architecture, PostgreSQL persistence, and a fully Dockerized environment.

## Features

- Create, retrieve, and track shipments via gRPC API
- Status event history with enforced state machine transitions
- PostgreSQL storage with versioned SQL migrations
- Manual dependency injection (no frameworks)
- Multi-stage Docker build with Docker Compose orchestration
- Graceful shutdown
- Unit and integration tests

## Tech Stack

- **Language:** Go 1.24
- **Transport:** gRPC (protobuf)
- **Database:** PostgreSQL 16 (via pgx)
- **Migrations:** golang-migrate
- **Containerization:** Docker, Docker Compose

## Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://docs.docker.com/get-docker/) and Docker Compose
- [protoc](https://grpc.io/docs/protoc-installation/) with Go plugins (only if regenerating proto)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (only for local migrations)

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/vektor-backend.git
cd vektor-backend
```

### 2. Configure environment

Create a `.env` file in the project root (or use the existing one):

```env
GRPC_PORT=50051

POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=shipments
POSTGRES_SSLMODE=disable
```

> When running the app **locally** (outside Docker), set `POSTGRES_HOST=127.0.0.1` and `POSTGRES_PORT=5434`.

### 3. Start with Docker Compose

```bash
make docker-up
```

This will:
1. Start PostgreSQL on port `5434` (host) / `5432` (container)
2. Run migrations automatically via the `migrator` service
3. Build and start the application on port `50051`

The service is ready when all three containers are healthy.

### 4. Stop everything

```bash
make docker-down
```

## Running Locally (without Docker for the app)

If you want to run only PostgreSQL in Docker and the app on your host:

```bash
# Start only Postgres
docker compose up -d postgres

# Apply migrations
make migrate-up

# Run the app
make run
```

## Project Structure

```text
cmd/
  shipment-service/
    main.go                  # Entrypoint

internal/
  app/
    app.go                   # Application bootstrap and lifecycle
    providers.go             # Dependency wiring

  config/
    config.go                # Environment-based configuration

  domain/
    shipment/
      entity.go              # Shipment aggregate
      event.go               # Status event value object
      status.go              # Status enum and transition rules
      errors.go              # Domain errors
      repository.go          # Repository interface

  repository/
    postgres/
      shipment_repository.go # PostgreSQL implementation

  usecase/
    shipment/
      create.go              # CreateShipment use case
      get.go                 # GetShipment use case
      add_event.go           # AddShipmentEvent use case
      get_history.go         # GetShipmentHistory use case

  transport/
    grpc/
      handler.go             # gRPC service implementation
      mapper.go              # Domain <-> protobuf mapping
      server.go              # gRPC server setup

proto/
  shipment.proto             # Service and message definitions
  shipment.pb.go             # Generated code
  shipment_grpc.pb.go        # Generated gRPC stubs

db/
  migrations/
    000001_init.up.sql       # Initial schema
    000001_init.down.sql     # Rollback

tests/
  integration/
    postgres_shipment_repository_test.go
```

## API Reference

The service exposes four gRPC methods on `localhost:50051`:

### CreateShipment

Creates a new shipment with initial `pending` status.

```json
{
  "reference": "REF-001",
  "origin": "Almaty",
  "destination": "Astana",
  "driver": "John Doe",
  "unit": "TRUCK-01",
  "amount": 1000,
  "driver_revenue": 700
}
```

### GetShipment

Returns a shipment by its reference.

```json
{
  "reference": "REF-001"
}
```

### AddShipmentEvent

Adds a status transition event. Invalid transitions are rejected.

```json
{
  "reference": "REF-001",
  "status": "SHIPMENT_STATUS_PICKED_UP"
}
```

### GetShipmentHistory

Returns the full status history for a shipment, ordered by time.

```json
{
  "reference": "REF-001"
}
```

## Status Transitions

Shipments follow a strict state machine. Only the transitions below are allowed:

| From        | Allowed next states      |
|-------------|--------------------------|
| pending     | picked_up, cancelled     |
| picked_up   | in_transit               |
| in_transit  | delivered                |
| delivered   | _(terminal state)_       |
| cancelled   | _(terminal state)_       |

Any other transition returns an error.

## Makefile Commands

| Command                | Description                        |
|------------------------|------------------------------------|
| `make proto`           | Regenerate gRPC code from `.proto` |
| `make migrate-up`      | Apply all pending migrations       |
| `make migrate-down`    | Roll back the last migration       |
| `make run`             | Run the app locally                |
| `make build`           | Build the project                  |
| `make test`            | Run unit tests                     |
| `make test-integration`| Run integration tests              |
| `make docker-up`       | Start all services via Docker      |
| `make docker-down`     | Stop and remove containers/volumes |
| `make docker-logs`     | Follow container logs              |

## Testing

### Unit tests

```bash
make test
```

### Integration tests

Integration tests require a running PostgreSQL instance:

```bash
docker compose up -d postgres
make test-integration
```

## Testing with Postman

Postman supports gRPC natively:

1. Create a new **gRPC Request**
2. Enter server address: `localhost:50051`
3. Import the proto file: `proto/shipment.proto`
4. Select a method and send your request
