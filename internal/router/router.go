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
		users.Delete("/:id", middleware.Authorize(model.RoleSuperAdmin, model.RoleCompanyRep, model.RoleAdmin), handlers.User.DeleteUser)
	}

	// Rute Super Admin
	superAdmin := authRequired.Group("/super-admin", middleware.Authorize(model.RoleSuperAdmin))
	{
		superAdmin.Post("/companies", handlers.User.CreateCompanyWithRep)
	}

	// Rute Company
	companies := authRequired.Group("/companies")
	{
		companies.Post("/", middleware.Authorize(model.RoleSuperAdmin), handlers.Company.CreateCompany)
		companies.Get("/", handlers.Company.GetAllCompanies)
		companies.Get("/:id", handlers.Company.GetCompany)
	}

	// Rute Store
	stores := authRequired.Group("/stores")
	{
		stores.Post("/", middleware.Authorize(model.RoleCompanyRep, model.RoleAdmin), handlers.Store.CreateStore)
		stores.Get("/", handlers.Store.GetAllStores)
		stores.Get("/:id", handlers.Store.GetStore)
		stores.Put("/:id", middleware.Authorize(model.RoleCompanyRep, model.RoleAdmin), handlers.Store.UpdateStore)
		stores.Delete("/:id", middleware.Authorize(model.RoleCompanyRep, model.RoleAdmin), handlers.Store.DeleteStore)
		stores.Post("/clone", middleware.Authorize(model.RoleCompanyRep, model.RoleAdmin), handlers.Store.CloneStoreContent)
	}

	// Rute Menu & Kategori
	menus := authRequired.Group("/menus", middleware.Authorize(model.RoleCompanyRep, model.RoleStoreAdmin, model.RoleAdmin))
	{
		menus.Post("/", handlers.Product.CreateMenu)
		menus.Get("/:id", handlers.Product.GetMenuByID)
		menus.Put("/:id", handlers.Product.UpdateMenu)
		menus.Delete("/:id", handlers.Product.DeleteMenu)
	}

	categories := authRequired.Group("/categories", middleware.Authorize(model.RoleCompanyRep, model.RoleStoreAdmin, model.RoleAdmin))
	{
		categories.Post("/", handlers.Product.CreateCategory)
		categories.Get("/:id", handlers.Product.GetCategoryByID)
		categories.Put("/:id", handlers.Product.UpdateCategory)
		categories.Delete("/:id", handlers.Product.DeleteCategory)
	}

	// Rute Promosi (CRUD Lengkap)
	promotions := authRequired.Group("/promotions", middleware.Authorize(model.RoleCompanyRep, model.RoleStoreAdmin, model.RoleAdmin))
	{
		// Discount Routes
		discounts := promotions.Group("/discounts")
		{
			discounts.Post("/", handlers.Promotion.CreateDiscount)
			discounts.Get("/", handlers.Promotion.ListDiscounts)
			discounts.Get("/:id", handlers.Promotion.GetDiscount)
			discounts.Put("/:id", handlers.Promotion.UpdateDiscount)
			discounts.Delete("/:id", handlers.Promotion.DeleteDiscount)
		}

		// Voucher Routes
		vouchers := promotions.Group("/vouchers")
		{
			vouchers.Post("/", handlers.Promotion.CreateVoucher)
			vouchers.Get("/", handlers.Promotion.ListVouchers)
			vouchers.Get("/:id", handlers.Promotion.GetVoucher)
			vouchers.Put("/:id", handlers.Promotion.UpdateVoucher)
			vouchers.Delete("/:id", handlers.Promotion.DeleteVoucher)
		}

		// Bundle Routes
		bundles := promotions.Group("/bundles")
		{
			bundles.Post("/", handlers.Promotion.CreateBundle)
			bundles.Get("/", handlers.Promotion.ListBundles)
			bundles.Get("/:id", handlers.Promotion.GetBundle)
			bundles.Put("/:id", handlers.Promotion.UpdateBundle)
			bundles.Delete("/:id", handlers.Promotion.DeleteBundle)
		}
	}

	// Rute Order
	orders := authRequired.Group("/orders")
	{
		orders.Post("/", handlers.Order.CreateOrder)
		orders.Get("/", handlers.Order.GetAllOrders)
		orders.Get("/:id", handlers.Order.GetOrder)
		orders.Post("/calculate-bill", handlers.Order.CalculateBill)
		orders.Patch("/:id/status", middleware.Authorize(model.RoleStoreAdmin, model.RoleCashier, model.RoleAdmin), handlers.Order.UpdateOrderStatus)
		orders.Put("/:id/items", middleware.Authorize(model.RoleStoreAdmin, model.RoleCashier, model.RoleAdmin), handlers.Order.UpdateOrderItems)
		orders.Post("/:id/confirm-payment", middleware.Authorize(model.RoleStoreAdmin, model.RoleCashier, model.RoleAdmin), handlers.Order.ConfirmPayment)
	}
}
