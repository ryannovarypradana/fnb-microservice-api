package main

import (
	"fnb-system/pkg/logger"
	"fnb-system/pkg/middleware"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	appLogger := logger.New()
	defer appLogger.Sync()

	app := fiber.New()
	app.Use(cors.New())

	// Tambahkan endpoint Health Check
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDR")
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	productServiceAddr := os.Getenv("PRODUCT_SERVICE_ADDR")
	orderServiceAddr := os.Getenv("ORDER_SERVICE_ADDR")

	api := app.Group("/api/v1")

	// Rute Publik
	api.Post("/register", proxy.Forward(authServiceAddr+"/api/v1/register"))
	api.Post("/login", proxy.Forward(authServiceAddr+"/api/v1/login"))
	api.Get("/menus", proxy.Forward(productServiceAddr+"/api/v1/menus"))
	api.Get("/categories", proxy.Forward(productServiceAddr+"/api/v1/categories"))

	// Rute Terproteksi
	api.Use(middleware.AuthMiddleware)

	api.Get("/users", proxy.Forward(userServiceAddr+"/api/v1/users"))
	api.Get("/users/:id", proxy.Forward(userServiceAddr+"/api/v1/users/:id"))
	api.Post("/orders", proxy.Forward(orderServiceAddr+"/api/v1/orders"))
	api.Get("/orders/my", proxy.Forward(orderServiceAddr+"/api/v1/orders/my"))
	api.Post("/menus", proxy.Forward(productServiceAddr+"/api/v1/menus"))
	api.Post("/categories", proxy.Forward(productServiceAddr+"/api/v1/categories"))

	appLogger.Info("API Gateway is starting on port 8080")
	if err := app.Listen(":8080"); err != nil {
		appLogger.Fatal("Failed to start API Gateway", zap.Error(err))
	}
}
