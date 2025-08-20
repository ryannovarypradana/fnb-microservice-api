// cmd/store-service/main.go

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	storepb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Get()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+cfg.Store.Port)
	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Store Svc on", cfg.Store.Port)

	companySvc, err := client.NewCompanyClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to company service: %v", err)
	}

	// PERBAIKAN: Menggunakan nama fungsi yang benar dari package 'internal/store'
	storeRepository := store.NewRepository(db)
	storeService := store.NewService(storeRepository, companySvc)
	storeGRPCHandler := store.NewGRPCHandler(storeService)

	grpcServer := grpc.NewServer()

	storepb.RegisterStoreServiceServer(grpcServer, storeGRPCHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
