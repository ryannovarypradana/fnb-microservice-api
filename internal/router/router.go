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

	api.Use(func(c *fiber.Ctx) error {
		// Periksa apakah method-nya adalah OPTIONS
		if c.Method() == "OPTIONS" {
			// Atur header yang sama persis dengan konfigurasi CORS Anda
			c.Set("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
			c.Set("Access-Control-Allow-Credentials", "true")
			// Kirim status 200 OK dan hentikan eksekusi lebih lanjut
			return c.SendStatus(fiber.StatusOK)
		}
		// Jika bukan OPTIONS, lanjutkan ke handler berikutnya
		return c.Next()
	})

	// Rute Publik
	auth := api.Group("/v1/auth")
	{
		auth.Post("/register", handlers.Auth.Register)
		auth.Post("/login", handlers.Auth.Login)
	}
	publicApi := api.Group("/public")
	{
		// Endpoint untuk Store
		publicApi.Get("/stores", handlers.Store.GetAllStores)
		publicApi.Get("/stores/:id", handlers.Store.GetStore)

		publicApi.Get("/stores/by-code/:storeCode", handlers.Store.GetStoreByCode)

		// Endpoint untuk Menu & Kategori berdasarkan Store
		publicApi.Get("/storse/by-code/:storeCode/menus", handlers.Product.GetMenusByStoreCode)
		publicApi.Get("/stores/by-code/:storeCode/categories", handlers.Product.GetCategoriesByStoreCode)

		// Endpoint untuk membuat pesanan publik
		publicApi.Post("/orders", handlers.Order.CreatePublicOrder)
	}

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

	pos := authRequired.Group("/pos", middleware.Authorize(model.RoleCashier, model.RoleStoreAdmin, model.RoleAdmin))
	{
		// Mengambil menu berdasarkan kode atau ID toko
		pos.Get("/stores/by-code/:storeCode/menus", handlers.Product.GetMenusByStoreCode)
		pos.Get("/stores/by-id/:id/menus", handlers.Product.GetMenusByStoreID)
		pos.Get("/menus/:id", handlers.Product.GetMenuByID)

		// Mengambil kategori berdasarkan kode atau ID toko
		pos.Get("/stores/by-code/:storeCode/categories", handlers.Product.GetCategoriesByStoreCode)
		pos.Get("/stores/by-id/:id/categories", handlers.Product.GetCategoriesByStoreID)
	}
}
