# Makefile untuk mengelola proyek microservices Go.

# Mengatur shell default
SHELL := /bin/bash

# --- DEVELOPMENT ---

## run-all: Menjalankan semua layanan di background dengan hot-reload (Air).
.PHONY: run-all
run-all:
	@echo "--- Starting all services with hot-reload ---"
	@echo "NOTE: Run 'make stop-all' in a separate terminal to stop."
	@trap 'make stop-all' EXIT
	@for service_path in $(wildcard cmd/*); do \
		SERVICE_NAME=$$(basename $$service_path); \
		echo "--> Starting $$SERVICE_NAME..."; \
		(export SERVICE_CMD_PATH=./$$service_path; air &) \
	done
	@echo "\n✅ All services are running in the background."
	@echo "   Use 'make stop-all' to terminate them."
	@echo "   Waiting for services to boot..."
	@# Keep the make command alive to manage background processes
	@tail -f /dev/null

## stop-all: Menghentikan semua proses Air dan binary yang berjalan.
.PHONY: stop-all
stop-all:
	@echo "--- Stopping all services ---"
	@pkill -f "air" > /dev/null 2>&1 || true
	@pkill -f "./tmp/main" > /dev/null 2>&1 || true
	@rm -rf ./tmp
	@echo "✅ All services stopped."

# --- gRPC CODE GENERATION ---

## proto-gen: Men-generate kode Go dari semua file .proto.
.PHONY: proto-gen
proto-gen:
	@echo "--- Generating gRPC code ---"
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/**/*.proto
	@echo "✅ gRPC code generated successfully."

# --- DOCKER ---

## docker-up: Menjalankan semua layanan menggunakan Docker Compose.
.PHONY: docker-up
docker-up:
	@echo "--- Starting all services with Docker Compose ---"
	@docker-compose up --build

## docker-down: Menghentikan semua layanan Docker Compose.
.PHONY: docker-down
docker-down:
	@echo "--- Stopping all services in Docker Compose ---"
	@docker-compose down

## docker-infra: Hanya menjalankan kontainer infrastruktur (DB, Redis, dll.).
.PHONY: docker-infra
docker-infra:
	@echo "--- Starting infrastructure containers ---"
	@docker-compose up -d postgres redis rabbitmq

## proto: Perintah alternatif untuk men-generate kode proto (opsional).
.PHONY: proto
proto:
	@echo "--- Cleaning and generating gRPC code ---"
	@rm -rf pkg/grpc/protoc
	@mkdir -p pkg/grpc/protoc
	@protoc -I=proto --go_out=. --go-grpc_out=. proto/**/*.proto