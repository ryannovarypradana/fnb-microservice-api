// /cmd/company-service/main.go
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	// FIX: Give the internal package a clear name, e.g., "internalCompany"
	internalCompany "github.com/ryannovarypradana/fnb-microservice-api/internal/company"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"

	// FIX: Give the gRPC package a clear alias, e.g., "companyPB" (for Protobuf)
	companyPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting Company Service...")

	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	if err := db.AutoMigrate(&model.Company{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed.")

	// Use the new package names
	companyRepo := internalCompany.NewRepository(db)
	companyService := internalCompany.NewService(companyRepo)
	companyHandler := internalCompany.NewGRPCHandler(companyService)

	port := os.Getenv("COMPANY_SERVICE_PORT")
	if port == "" {
		port = "50054"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// FIX: Use the alias for the gRPC package to register the server
	companyPB.RegisterCompanyServiceServer(grpcServer, companyHandler)
	reflection.Register(grpcServer)

	log.Printf("Company gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
