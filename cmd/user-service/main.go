package main

import (
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
	userSvc := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userSvc)

	app := fiber.New()
	app.Use(cors.New())

	// Tambahkan endpoint Health Check
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "user-service",
		})
	})

	api := app.Group("/api/v1")
	userRoutes := api.Group("/users")

	userRoutes.Get("", userHandler.GetAllUsers)
	userRoutes.Get("/:id", userHandler.GetUserByID)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Delete("/:id", userHandler.DeleteUser)

	appLogger.Info("User Service is starting on port 8082")
	if err := app.Listen(":8082"); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
