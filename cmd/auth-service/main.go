package main

import (
	"fnb-system/internal/auth"
	"fnb-system/internal/user"
	"fnb-system/pkg/database"
	"fnb-system/pkg/logger"
	"log"

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

	userRepo := user.NewUserRepository(db)
	authRepo := auth.NewAuthRepository(db)
	authSvc := auth.NewAuthService(authRepo, userRepo)
	authHandler := auth.NewAuthHandler(authSvc)

	app := fiber.New()
	app.Use(cors.New())

	// Tambahkan endpoint Health Check
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	api := app.Group("/api/v1")
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	appLogger.Info("Auth Service is starting on port 8081")
	if err := app.Listen(":8081"); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
