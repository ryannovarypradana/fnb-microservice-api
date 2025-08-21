# Makefile untuk mengelola proyek microservices Go.

# Mengatur shell default
SHELL := /bin/bash

# --- DEVELOPMENT ---

## run-all: Menjalankan semua layanan di background dengan hot-reload (Air).
.PHONY: run-all
run-all:
	@echo "--- Starting all services with hot-reload ---"
	@for service_path in $(wildcard cmd/*); do \
		SERVICE_NAME=$$(basename $$service_path); \
		echo "--> Starting $$SERVICE_NAME..."; \
		(export SERVICE_CMD_PATH=./$$service_path; air &) \
	done
	@echo "\n✅ All services are running in the background."
	@echo "   Use 'make stop-all' to terminate them."
	@echo "   Waiting for services to boot..."
	@sleep 2

## stop-all: Menghentikan semua proses Air dan binary yang berjalan.
.PHONY: stop-all
stop-all:
	@echo "--- Stopping all services ---"
	@killall -q air || true
	@killall -q main || true
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
	docker-compose up

## docker-down: Menghentikan semua layanan Docker Compose.
.PHONY: docker-down
docker-down:
	@echo "--- Stopping all services in Docker Compose ---"
	docker-compose down

## docker-infra: Hanya menjalankan kontainer infrastruktur (DB, Redis, dll.).
.PHONY: docker-infra
docker-infra:
	@echo "--- Starting infrastructure containers ---"
	docker-compose up -d postgres redis rabbitmq