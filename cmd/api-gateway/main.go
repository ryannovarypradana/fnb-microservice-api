// cmd/api-gateway/main.go

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/internal/router"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/handler"
)

func main() {
	// 1. INISIALISASI
	cfg := config.Get()

	authSvc, err := client.NewAuthClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to auth service: %v", err)
	}

	userSvc, err := client.NewUserClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to user service: %v", err)
	}

	storeSvc, err := client.NewStoreClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to store service: %v", err)
	}

	productSvc, err := client.NewProductClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to product service: %v", err)
	}

	orderSvc, err := client.NewOrderClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to order service: %v", err)
	}

	companySvc, err := client.NewCompanyClient(cfg) // <-- Inisialisasi Company Client
	if err != nil {
		log.Fatalf("Could not connect to company service: %v", err)
	}
	promotionSvc, err := client.NewPromotionClient(cfg)
	if err != nil {
		log.Fatalf("Could not connect to promotion service: %v", err)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Izinkan frontend Anda
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	// Kumpulkan semua handler
	handlers := &handler.Handlers{
		Auth:  handler.NewAuthHandler(authSvc),
		User:  handler.NewUserHandler(userSvc),
		Store: handler.NewStoreHandler(storeSvc),
		// DIperbaiki: Menambahkan storeSvc sebagai argumen kedua
		Product:   handler.NewProductHandler(productSvc, storeSvc),
		Order:     handler.NewOrderHandler(orderSvc),
		Company:   handler.NewCompanyHandler(companySvc), // <-- Tambahkan Company Handler
		Promotion: handler.NewPromotionHandler(promotionSvc),
	}

	// 2. SETUP ROUTE
	router.SetupRoutes(app, cfg, handlers)

	// 3. JALANKAN SERVER
	log.Printf("Starting API Gateway on port %s", cfg.App.Port)
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
