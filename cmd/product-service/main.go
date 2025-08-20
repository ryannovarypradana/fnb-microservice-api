package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	// Perbaikan di sini: Menambahkan semua import yang diperlukan
	"github.com/ryannovarypradana/fnb-microservice-api/config"
	internalProduct "github.com/ryannovarypradana/fnb-microservice-api/internal/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/database"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/logger"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	cfg := config.Get()
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Menggunakan model.Menu dan model.Category
	db.AutoMigrate(&model.Menu{}, &model.Category{}, &model.Store{})
	logger.InitLogger()

	lis, err := net.Listen("tcp", ":"+cfg.Product.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	productRepository := internalProduct.NewProductRepository(db)
	productService := internalProduct.NewProductService(productRepository)
	productGRPCHandler := internalProduct.NewProductGRPCHandler(productService)

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, productGRPCHandler)

	log.Printf("Product service is running on port %s", cfg.Product.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
