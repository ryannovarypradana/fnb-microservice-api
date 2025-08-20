package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	// CORRECTED: Import the single 'auth' package
	"github.com/ryannovarypradana/fnb-microservice-api/internal/auth"
	// CORRECTED: Fixed typo 'github.comcom'
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	// CORRECTED: Added missing import
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/utils"
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Auth.Port))
	if err != nil {
		log.Fatalf("FATAL: failed to listen on port %s: %v", cfg.Auth.Port, err)
	}

	jwtService := utils.NewJwtService(cfg)
	grpcServer := grpc.NewServer()

	// CORRECTED: All constructors are now called from the 'auth' package
	authRepo := auth.NewAuthRepository(db)
	authService := auth.NewAuthService(authRepo, jwtService)
	auth.NewAuthGRPCHandler(grpcServer, authService)

	log.Printf("Auth service is listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("FATAL: failed to serve gRPC server: %v", err)
	}
}
