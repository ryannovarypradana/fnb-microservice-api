// /cmd/auth-service/main.go
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/auth"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting Auth Service...")

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

	// Migrasi User model. Pastikan semua relasi sudah didefinisikan di pkg/model
	if err := db.AutoMigrate(&model.User{}, &model.Company{}, &model.Store{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed.")

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, db)
	authHandler := auth.NewGRPCHandler(authService)

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "50051" // Port default
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)
	reflection.Register(grpcServer)

	log.Printf("Auth gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
