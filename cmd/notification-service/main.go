// cmd/notification-service/main.go
package main

import (
	"fnb-system/internal/notification"
	"os"
)

func main() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	// Jalankan consumer
	notification.StartOrderCreatedConsumer(rabbitMQURL)
}
