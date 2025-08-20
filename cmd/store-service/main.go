package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/ryannovarypradana/fnb-microservice-api/config"
	internalStore "github.com/ryannovarypradana/fnb-microservice-api/internal/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client" // <-- Import client
	storePB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting Store Service...")

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	cfg := config.Get()
	// Perbaikan di sini: Teruskan cfg ke NewPostgres
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	if err := db.AutoMigrate(&model.Store{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// 1. Buat koneksi gRPC ke company-service
	companyClient := client.NewCompanyClient()

	// 2. Inisialisasi service dengan menyertakan companyClient
	storeRepo := internalStore.NewRepository(db)
	storeService := internalStore.NewService(storeRepo, companyClient) // <-- Suntikkan di sini
	storeHandler := internalStore.NewGRPCHandler(storeService)

	port := os.Getenv("STORE_SERVICE_PORT")
	if port == "" {
		port = "50055"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	storePB.RegisterStoreServiceServer(grpcServer, storeHandler)
	reflection.Register(grpcServer)

	log.Printf("Store gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
