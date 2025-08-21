package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/handler"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/middleware"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, handlers *handler.Handlers) {
	api := app.Group("/api")

	// Rute Publik
	auth := api.Group("/auth")
	{
		auth.Post("/register", handlers.Auth.Register)
		auth.Post("/login", handlers.Auth.Login)
	}
	api.Post("/public/orders", handlers.Order.CreatePublicOrder)

	// --- Rute Terproteksi ---
	authRequired := api.Group("/v1", middleware.AuthMiddleware(cfg))

	// Rute User
	users := authRequired.Group("/users")
	{
		users.Get("/:id", handlers.User.GetUser)
		users.Put("/:id", handlers.User.UpdateUser)
		// User dihapus oleh peran administratif
		users.Delete("/:id", middleware.Authorize(model.RoleSuperAdmin, model.RoleCompanyRep, model.RoleAdmin), handlers.User.DeleteUser)
	}

	// Rute Super Admin (Hanya SUPER_ADMIN)
	superAdmin := authRequired.Group("/super-admin", middleware.Authorize(model.RoleSuperAdmin))
	{
		superAdmin.Post("/companies", handlers.User.CreateCompanyWithRep)
	}

	// Rute Company
	companies := authRequired.Group("/companies")
	{
		// Hanya SUPER_ADMIN yang bisa membuat company
		companies.Post("/", middleware.Authorize(model.RoleSuperAdmin), handlers.Company.CreateCompany)
		companies.Get("/", handlers.Company.GetAllCompanies)
		companies.Get("/:id", handlers.Company.GetCompany)
	}

	// Rute Store (Dikelola oleh COMPANY_REP dan ADMIN)
	stores := authRequired.Group("/stores")
	{
		stores.Post("/", middleware.Authorize(model.RoleCompanyRep, model.RoleAdmin), handlers.Store.CreateStore)
		stores.Get("/", handlers.Store.GetAllStores)
		stores.Get("/:id", handlers.Store.GetStore)
		stores.Put("/:id", middleware.Authorize(model.RoleCompanyRep, model.RoleAdmin), handlers.Store.UpdateStore)
		stores.Delete("/:id", middleware.Authorize(model.RoleCompanyRep, model.RoleAdmin), handlers.Store.DeleteStore)
		stores.Post("/clone", middleware.Authorize(model.RoleCompanyRep, model.RoleAdmin), handlers.Store.CloneStoreContent)
	}

	// Rute Menu (Dikelola oleh COMPANY_REP, STORE_ADMIN, dan ADMIN)
	menus := authRequired.Group("/menus")
	{
		menus.Post("/", middleware.Authorize(model.RoleCompanyRep, model.RoleStoreAdmin, model.RoleAdmin), handlers.Product.CreateMenu)
		menus.Get("/:id", handlers.Product.GetMenuByID)
	}

	// Rute Order
	orders := authRequired.Group("/orders")
	{
		// Semua pengguna yang login dapat membuat pesanan
		orders.Post("/", handlers.Order.CreateOrder)
		orders.Get("/", handlers.Order.GetAllOrders)
		orders.Get("/:id", handlers.Order.GetOrder)
		orders.Post("/calculate-bill", handlers.Order.CalculateBill)
		// Operasional order oleh STORE_ADMIN, CASHIER, dan ADMIN
		orders.Patch("/:id/status", middleware.Authorize(model.RoleStoreAdmin, model.RoleCashier, model.RoleAdmin), handlers.Order.UpdateOrderStatus)
		orders.Put("/:id/items", middleware.Authorize(model.RoleStoreAdmin, model.RoleCashier, model.RoleAdmin), handlers.Order.UpdateOrderItems)
		orders.Post("/:id/confirm-payment", middleware.Authorize(model.RoleStoreAdmin, model.RoleCashier, model.RoleAdmin), handlers.Order.ConfirmPayment)
	}
}
