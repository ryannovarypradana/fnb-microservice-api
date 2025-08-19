package main

import (
	"fmt"
	"log"
	"net"
	"os"

	internalOrder "github.com/ryannovarypradana/fnb-microservice-api/internal/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/eventbus" // <-- Import eventbus
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	orderPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting Order Service...")

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.Order{}, &model.OrderItem{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	productClient := client.NewProductClient()

	// Inisialisasi RabbitMQ Publisher
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	publisher, err := eventbus.NewRabbitMQPublisher(rabbitMQURL)
	if err != nil {
		log.Fatalf("gagal terhubung ke RabbitMQ: %v", err)
	}
	defer publisher.Close()

	// Suntikkan semua dependensi
	orderRepo := internalOrder.NewRepository(db)
	orderService := internalOrder.NewService(orderRepo, productClient, publisher) // <-- Diperbarui
	orderHandler := internalOrder.NewGRPCHandler(orderService)

	port := os.Getenv("ORDER_SERVICE_PORT")
	if port == "" {
		port = "50056"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderPB.RegisterOrderServiceServer(grpcServer, orderHandler)
	reflection.Register(grpcServer)

	log.Printf("Order gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
