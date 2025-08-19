// /cmd/api-gateway/main.go
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/handler"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/middleware"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found.")
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// === Initialize gRPC Clients ===
	authClient := client.NewAuthClient()
	userClient := client.NewUserClient()
	// ... other clients

	// === Initialize Handlers ===
	authHandler := handler.NewAuthHandler(authClient)
	// FIX: Initialize the userHandler using the userClient
	userHandler := handler.NewUserHandler(userClient)
	// ... other handlers

	// === Setup Routes ===
	api := app.Group("/api")

	// --- Public Routes ---
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)

	// --- Authenticated Routes ---
	authRequired := api.Group("", middleware.AuthMiddleware)

	// FIX: Use the initialized userHandler and the correct Role constant
	authRequired.Get("/users/:id", middleware.RoleMiddleware(model.RoleSuperAdmin), userHandler.GetUser)

	// === Start Server ===
	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
