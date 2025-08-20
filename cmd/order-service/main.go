package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/logger"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	cfg := config.Get()
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	db.AutoMigrate(&model.Order{}, &model.OrderItem{})
	logger.InitLogger()

	lis, err := net.Listen("tcp", ":"+cfg.Order.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	productClient := client.NewProductServiceClient(cfg.Product.Port)
	orderRepository := order.NewOrderRepository(db)
	orderService := order.NewOrderService(orderRepository, productClient)
	orderGRPCHandler := order.NewOrderGRPCHandler(orderService)

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderGRPCHandler)

	log.Printf("Order service is running on port %s", cfg.Order.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
