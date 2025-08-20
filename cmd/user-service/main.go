// cmd/user-service/main.go

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/user"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/eventbus"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	userpb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Get()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// PERBAIKAN: Inisialisasi RabbitMQ dengan cara yang benar
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.Rabbit.User,
		cfg.Rabbit.Password,
		cfg.Rabbit.Host,
		cfg.Rabbit.Port,
	)
	publisher, err := eventbus.NewRabbitMQPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer publisher.Close()

	lis, err := net.Listen("tcp", ":"+cfg.User.Port)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("User Svc on", cfg.User.Port)

	storeSvc, err := client.NewStoreClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to store service: %v", err)
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository, storeSvc, publisher)
	userGRPCHandler := user.NewGRPCHandler(userService)

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, userGRPCHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
