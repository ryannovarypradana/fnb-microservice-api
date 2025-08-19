package main

import (
	"fmt"
	"log"
	"net"
	"os"

	internalProduct "github.com/ryannovarypradana/fnb-microservice-api/internal/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client" // <-- Import paket client
	productPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting Product Service...")

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.Product{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 1. Inisialisasi gRPC client untuk store-service
	storeClient := client.NewStoreClient()

	// 2. Suntikkan (inject) storeClient ke dalam productService
	productRepo := internalProduct.NewRepository(db)
	productService := internalProduct.NewService(productRepo, storeClient) // <-- Diperbarui
	productHandler := internalProduct.NewGRPCHandler(productService)

	port := os.Getenv("PRODUCT_SERVICE_PORT")
	if port == "" {
		port = "50053"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	productPB.RegisterProductServiceServer(grpcServer, productHandler)
	reflection.Register(grpcServer)

	log.Printf("Product gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
