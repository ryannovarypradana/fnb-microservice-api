package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/notification"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found.")
	}

	log.Println("Starting Notification Service...")

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatalf("RABBITMQ_URL environment variable not set")
	}

	consumer, err := notification.NewRabbitMQConsumer(rabbitMQURL)
	if err != nil {
		log.Fatalf("Gagal membuat consumer RabbitMQ: %v", err)
	}

	// StartConsuming akan berjalan selamanya (blocking)
	consumer.StartConsuming()
}
