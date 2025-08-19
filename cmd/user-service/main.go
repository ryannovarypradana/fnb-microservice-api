package main

import (
	"fmt"
	"log"
	"net"
	"os"

	internalUser "github.com/ryannovarypradana/fnb-microservice-api/internal/user"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"

	// FIX: Correctly import the eventbus package
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/eventbus"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	userPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting User Service...")

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Initialize gRPC client for other services
	storeClient := client.NewStoreClient()

	// Initialize RabbitMQ Publisher
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	// FIX: Use the correct package name to call the function
	publisher, err := eventbus.NewRabbitMQPublisher(rabbitMQURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer publisher.Close()

	// Inject all dependencies into the service
	userRepo := internalUser.NewRepository(db)
	userService := internalUser.NewService(userRepo, storeClient, publisher)
	userHandler := internalUser.NewGRPCHandler(userService)

	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userPB.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	log.Printf("User gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
