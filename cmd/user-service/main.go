// cmd/user-service/main.go
package main

import (
	"fmt"
	"log"
	"net"
	"os" // <-- TAMBAHKAN IMPORT

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/user"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	companyPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
)

func main() {
	cfg := config.Get()
	if cfg == nil {
		log.Fatalf("FATAL: could not load config")
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("FATAL: failed to connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.User.Port))
	if err != nil {
		log.Fatalf("FATAL: failed to listen on port %s: %v", cfg.User.Port, err)
	}

	// --- PERUBAHAN DI SINI ---
	companyServiceAddr := os.Getenv("COMPANY_SERVICE_ADDR")
	if companyServiceAddr == "" {
		log.Fatalf("FATAL: COMPANY_SERVICE_ADDR environment variable not set")
	}

	companyConn, err := grpc.Dial(
		companyServiceAddr, // <-- Gunakan variabel lingkungan
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("FATAL: failed to connect to company service: %v", err)
	}
	defer companyConn.Close()
	companyClient := companyPB.NewCompanyServiceClient(companyConn)
	log.Println("Successfully connected to Company service.")
	// --- AKHIR PERUBAHAN ---

	grpcServer := grpc.NewServer()

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo, companyClient)
	user.NewUserGRPCHandler(grpcServer, userService)

	log.Printf("User service is listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("FATAL: failed to serve gRPC server: %v", err)
	}
}
