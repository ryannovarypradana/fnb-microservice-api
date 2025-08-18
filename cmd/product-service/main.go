package main

import (
	"fnb-system/internal/product"
	"fnb-system/pkg/database"
	"fnb-system/pkg/logger"
	"fnb-system/pkg/redis"
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewRedisClient(redisAddr)

	productRepo := product.NewProductRepository(db, rdb)
	productSvc := product.NewProductService(productRepo)
	productHandler := product.NewProductHandler(productSvc)

	app := fiber.New()
	app.Use(cors.New())

	// Tambahkan endpoint Health Check
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "product-service",
		})
	})

	api := app.Group("/api/v1")

	api.Get("/categories", productHandler.GetAllCategories)
	api.Get("/menus", productHandler.GetAllMenus)
	api.Get("/menus/:id", productHandler.GetMenuByID)

	adminRoutes := api.Group("/")
	adminRoutes.Post("/categories", productHandler.CreateCategory)
	adminRoutes.Post("/menus", productHandler.CreateMenu)

	appLogger.Info("Product Service is starting on port 8083")
	if err := app.Listen(":8083"); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
