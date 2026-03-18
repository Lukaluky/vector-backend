PROTO_DIR=proto
PROTO_FILE=$(PROTO_DIR)/shipment.proto
DB_URL=postgres://postgres:postgres@localhost:5432/shipments?sslmode=disable

proto:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILE)

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down 1

run:
	go run ./cmd/shipment-service

build:
	go build ./...