// cmd/order-service/main.go
package main

import (
	"fmt"
	"log"
	"net"
	"os" // <-- TAMBAHKAN IMPORT

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	productPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	storePB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
)

// Helper function untuk koneksi gRPC yang lebih bersih
func dialService(envVarName string) *grpc.ClientConn {
	addr := os.Getenv(envVarName)
	if addr == "" {
		log.Fatalf("FATAL: environment variable %s not set", envVarName)
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("FATAL: failed to connect to %s: %v", addr, err)
	}
	return conn
}

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

	// --- PERUBAHAN DI SINI ---
	// Membuat koneksi klien gRPC ke service lain menggunakan helper
	productConn := dialService("PRODUCT_SERVICE_ADDR")
	defer productConn.Close()
	productClient := productPB.NewProductServiceClient(productConn)
	log.Println("Successfully connected to Product service.")

	storeConn := dialService("STORE_SERVICE_ADDR")
	defer storeConn.Close()
	storeClient := storePB.NewStoreServiceClient(storeConn)
	log.Println("Successfully connected to Store service.")
	// --- AKHIR PERUBAHAN ---

	// Menyiapkan server gRPC untuk service ini
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Order.Port))
	if err != nil {
		log.Fatalf("FATAL: failed to listen on port %s: %v", cfg.Order.Port, err)
	}

	grpcServer := grpc.NewServer()

	// Inisialisasi dari paket 'order'
	orderRepo := order.NewOrderRepository(db)
	orderService := order.NewOrderService(orderRepo, db, productClient, storeClient)
	order.NewOrderGRPCHandler(grpcServer, orderService)

	log.Printf("Order service is listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("FATAL: failed to serve gRPC server: %v", err)
	}
}
