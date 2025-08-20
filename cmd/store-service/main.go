package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
)

func main() {
	// Memuat konfigurasi dari file .env
	// CORRECTED: Called config.Get() instead of config.LoadConfig()
	cfg := config.Get()
	if cfg == nil {
		log.Fatalf("FATAL: could not load config")
	}

	// Menghubungkan ke database PostgreSQL
	// CORRECTED: Called database.NewPostgres() instead of database.NewPostgresDB()
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("FATAL: failed to connect to database: %v", err)
	}
	log.Println("Database connection successful.")

	// Membuat listener untuk koneksi gRPC
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Store.Port))
	if err != nil {
		log.Fatalf("FATAL: failed to listen on port %s: %v", cfg.Store.Port, err)
	}

	// Membuat server gRPC baru
	grpcServer := grpc.NewServer()

	// Inisialisasi dan wiring dependensi dari paket 'store'
	// 1. Buat repository
	storeRepo := store.NewStoreRepository(db)
	// 2. Buat service dengan repository
	storeService := store.NewStoreService(storeRepo)
	// 3. Daftarkan handler gRPC dengan service
	store.NewStoreGRPCHandler(grpcServer, storeService)

	log.Printf("Store service is listening at %v", lis.Addr())

	// Mulai melayani request gRPC
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("FATAL: failed to serve gRPC server: %v", err)
	}
}
