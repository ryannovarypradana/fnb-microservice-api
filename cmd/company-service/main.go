// cmd/company-service/main.go
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	internalCompany "github.com/ryannovarypradana/fnb-microservice-api/internal/company"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	companyPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Starting Company Service...")

	// Memuat konfigurasi terlebih dahulu
	cfg := config.Get()
	if cfg == nil {
		log.Fatalf("FATAL: Config object could not be loaded.")
	}

	// === BARIS DEBUGGING ===
	// Baris ini akan mencetak HOST database yang akan digunakan.
	// Perhatikan baik-baik output dari baris ini saat Anda menjalankan service.
	log.Printf("DEBUG: Attempting to connect to DB_HOST: '%s'", cfg.DB.Host)
	// =======================

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	if err := db.AutoMigrate(&model.Company{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed.")

	companyRepo := internalCompany.NewRepository(db)
	companyService := internalCompany.NewService(companyRepo)
	companyHandler := internalCompany.NewGRPCHandler(companyService)

	port := os.Getenv("COMPANY_SERVICE_PORT")
	if port == "" {
		port = "50053" // Fallback port
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	companyPB.RegisterCompanyServiceServer(grpcServer, companyHandler)
	reflection.Register(grpcServer)

	log.Printf("Company gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}
