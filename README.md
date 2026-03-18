# Shipment Service

gRPC service for shipment management with PostgreSQL persistence and migrations.

## Features

- Create shipment
- Get shipment by reference
- Add shipment status event
- Get shipment history
- PostgreSQL storage
- SQL migrations
- Graceful shutdown
- Manual dependency injection
- Clean architecture style structure

## Tech Stack

- Go
- gRPC
- PostgreSQL
- pgx
- golang-migrate
- Docker Compose

## Project Structure

```text
cmd/
  shipment-service/
    main.go

internal/
  app/
    app.go
    providers.go

  config/
    config.go

  domain/
    shipment/
      entity.go
      event.go
      status.go
      errors.go
      repository.go

  repository/
    postgres/
      shipment_repository.go

  usecase/
    shipment/
      create.go
      get.go
      add_event.go
      get_history.go

  transport/
    grpc/
      handler.go
      mapper.go
      server.go

proto/
  shipment.proto
  shipment.pb.go
  shipment_grpc.pb.go

db/
  migrations/
    000001_init.up.sql
    000001_init.down.sql