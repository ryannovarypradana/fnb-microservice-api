package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	internalProduct "github.com/ryannovarypradana/fnb-microservice-api/internal/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/logger"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

func main() {
	log.Println("Starting Product Service...")

	// 1. Muat Konfigurasi
	// Memanggil config.Get() untuk memuat semua variabel lingkungan dari file .env
	cfg := config.Get()
	if cfg == nil {
		log.Fatalf("FATAL: Config object could not be loaded.")
	}

	// 2. Hubungkan ke Database
	// Membuat koneksi ke PostgreSQL menggunakan konfigurasi yang sudah dimuat.
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Database connection successful.")

	// 3. Lakukan Migrasi Database
	// AutoMigrate akan membuat atau memperbarui tabel sesuai dengan struct Go.
	db.AutoMigrate(&model.Menu{}, &model.Category{}, &model.Store{})

	// Inisialisasi logger jika diperlukan
	logger.InitLogger()

	// 4. Siapkan Listener untuk gRPC
	// Service ini akan "mendengarkan" permintaan gRPC yang masuk pada port yang ditentukan.
	lis, err := net.Listen("tcp", ":"+cfg.Product.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Product.Port, err)
	}

	// 5. Inisialisasi Semua Dependensi (Repository, Service, Handler)
	// Ini adalah inti dari Dependency Injection: membuat komponen dari yang paling dalam (repository) ke yang paling luar (handler).
	productRepository := internalProduct.NewProductRepository(db)
	productService := internalProduct.NewProductService(productRepository)
	productGRPCHandler := internalProduct.NewProductGRPCHandler(productService)

	// 6. Buat dan Daftarkan Server gRPC
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, productGRPCHandler)

	// 7. Jalankan Server
	// Server akan mulai berjalan dan siap menerima koneksi.
	log.Printf("Product service is running on port %s", cfg.Product.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
