package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/handler"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/middleware"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, handlers *handler.Handlers) {
	api := app.Group("/api")

	// Rute Auth
	api.Post("/auth/register", handlers.Auth.Register)
	api.Post("/auth/login", handlers.Auth.Login)

	// --- Rute yang Membutuhkan Autentikasi ---
	authRequired := api.Group("/v1", middleware.AuthMiddleware(cfg))

	// Rute untuk User
	users := authRequired.Group("/users")
	{
		users.Get("/:id", handlers.User.GetUser)
		// Memerlukan otorisasi, misalnya hanya admin atau user itu sendiri
		users.Put("/:id", handlers.User.UpdateUser)
		// Memerlukan otorisasi, misalnya hanya admin
		users.Delete("/:id", handlers.User.DeleteUser)
	}

	// Rute untuk Super Admin
	superAdmin := authRequired.Group("/super-admin") // Tambahkan middleware otorisasi peran di sini
	{
		superAdmin.Post("/companies", handlers.User.CreateCompanyWithRep)
	}

	// Rute untuk Menu
	authRequired.Post("/menus", handlers.Product.CreateMenu)
	authRequired.Get("/menus/:id", handlers.Product.GetMenuByID)

	// Rute untuk Order
	authRequired.Post("/orders", handlers.Order.CreateOrder)

	// Rute untuk Store & Company
	authRequired.Post("/stores", handlers.Store.CreateStore)
	authRequired.Post("/companies", handlers.Company.CreateCompany)
	authRequired.Get("/companies/:id", handlers.Company.GetCompany)

}
