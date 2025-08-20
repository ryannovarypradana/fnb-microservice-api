// cmd/order-service/main.go

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	orderpb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Get()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+cfg.Order.Port)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Order Svc on", cfg.Order.Port)

	productSvc, err := client.NewProductClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to product service: %v", err)
	}

	orderRepository := order.NewOrderRepository(db)
	orderService := order.NewOrderService(orderRepository, productSvc)

	// DI SINILAH PEMANGGILAN ITU TERJADI
	orderGRPCHandler := order.NewOrderGRPCHandler(orderService)

	grpcServer := grpc.NewServer()

	orderpb.RegisterOrderServiceServer(grpcServer, orderGRPCHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
