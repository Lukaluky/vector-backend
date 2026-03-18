PROTO_DIR=proto
PROTO_FILE=$(PROTO_DIR)/shipment.proto
DB_URL=postgres://postgres:postgres@127.0.0.1:5434/shipments?sslmode=disable

proto:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILE)

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down 1

run:
	go run ./cmd/shipment-service

build:
	go mod tidy
	go build ./...

test:
	go test ./...

test-integration:
	go test ./tests/integration -v

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down -v

docker-logs:
	docker compose logs -f