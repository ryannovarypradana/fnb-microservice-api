package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/promotion"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	promotionpb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/promotion"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc"
)

func main() {
	// 1. Load Konfigurasi dari .env
	cfg := config.Get()
	if cfg == nil {
		log.Fatalf("FATAL: Config object could not be loaded.")
	}

	// 2. Koneksi ke Database
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Database connection successful.")

	// 3. Auto-Migrate Model Database
	// GORM akan membuat tabel discounts, vouchers, dan bundles jika belum ada.
	log.Println("Running database migrations...")
	err = db.AutoMigrate(&model.Discount{}, &model.Voucher{}, &model.Bundle{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 4. Inisialisasi TCP Listener
	// Service akan mendengarkan koneksi gRPC pada alamat yang ditentukan di .env

	lis, err := net.Listen("tcp", ":"+cfg.Promotion.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Promotion.Port, err)
	}

	// 5. Inisialisasi Repository, Service, dan Handler
	// Ini adalah dependency injection, menyambungkan semua lapisan aplikasi.
	promoRepository := promotion.NewRepository(db)
	promoService := promotion.NewService(promoRepository)
	promoHandler := promotion.NewGrpcHandler(promoService)

	// 6. Membuat dan Mendaftarkan gRPC Server
	grpcServer := grpc.NewServer()
	promotionpb.RegisterPromotionServiceServer(grpcServer, promoHandler)

	// 7. Menjalankan Server
	fmt.Printf("✅ Promotion service is running and listening on %s\n", cfg.Promotion.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}

}
