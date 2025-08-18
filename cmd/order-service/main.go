package main

import (
	"fnb-system/internal/order"
	"fnb-system/pkg/database"
	"fnb-system/pkg/eventbus"
	"fnb-system/pkg/logger"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	appLogger := logger.New()
	defer appLogger.Sync()

	db, err := database.NewPostgresConnection()
	if err != nil {
		appLogger.Fatal("Failed to connect to database", zap.Error(err))
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}
	bus, err := eventbus.NewRabbitMQBus(rabbitMQURL)
	if err != nil {
		appLogger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}
	defer bus.Close()

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://localhost:8083"
	}

	productClient := order.NewProductClient(productServiceURL)
	orderRepo := order.NewOrderRepository(db)
	orderSvc := order.NewOrderService(orderRepo, productClient, bus)
	orderHandler := order.NewOrderHandler(orderSvc)

	app := fiber.New()
	app.Use(cors.New())

	// Tambahkan endpoint Health Check
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "order-service",
		})
	})

	api := app.Group("/api/v1")

	orderRoutes := api.Group("/orders")
	orderRoutes.Post("", orderHandler.CreateOrder)
	orderRoutes.Get("/my", orderHandler.GetMyOrders)

	appLogger.Info("Order Service is starting on port 8084")
	if err := app.Listen(":8084"); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
