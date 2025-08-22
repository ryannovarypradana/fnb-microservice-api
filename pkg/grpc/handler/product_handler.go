package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid" // <-- Pastikan import ini ada
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/dto"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
)

// ProductHandler adalah interface untuk handler produk (Menu dan Kategori) di API Gateway.
type ProductHandler interface {
	// Menu Handlers
	CreateMenu(c *fiber.Ctx) error
	GetMenuByID(c *fiber.Ctx) error
	UpdateMenu(c *fiber.Ctx) error
	DeleteMenu(c *fiber.Ctx) error

	// Category Handlers
	CreateCategory(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error
	UpdateCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error

	// --- POS Handlers ---
	GetMenusByStoreCode(c *fiber.Ctx) error
	GetMenusByStoreID(c *fiber.Ctx) error
	GetCategoriesByStoreCode(c *fiber.Ctx) error
	GetCategoriesByStoreID(c *fiber.Ctx) error
}

type productHandler struct {
	productClient product.ProductServiceClient
	storeClient   store.StoreServiceClient
}

// NewProductHandler membuat instance baru dari productHandler.
func NewProductHandler(productClient product.ProductServiceClient, storeClient store.StoreServiceClient) ProductHandler {
	return &productHandler{
		productClient: productClient,
		storeClient:   storeClient,
	}
}

// ============== MENU HANDLERS ==============

// CreateMenu menangani permintaan HTTP untuk membuat menu baru.
func (h *productHandler) CreateMenu(c *fiber.Ctx) error {
	var req dto.CreateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	// Validasi bahwa StoreID ada di body dan merupakan UUID yang valid.
	if _, err := uuid.Parse(req.StoreID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or missing storeId in request body",
		})
	}

	// Validasi CategoryID jika diperlukan
	if _, err := uuid.Parse(req.CategoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or missing categoryId in request body",
		})
	}

	grpcRequest := &product.CreateMenuRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageUrl:    req.ImageURL,
		CategoryId:  req.CategoryID,
		StoreId:     req.StoreID,
	}

	res, err := h.productClient.CreateMenu(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetMenuByID menangani permintaan HTTP untuk mencari satu menu berdasarkan ID.
func (h *productHandler) GetMenuByID(c *fiber.Ctx) error {
	menuID := c.Params("id")

	grpcRequest := &product.GetMenuByIDRequest{
		MenuId: menuID,
	}

	res, err := h.productClient.GetMenuByID(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateMenu menangani permintaan HTTP untuk memperbarui menu.
func (h *productHandler) UpdateMenu(c *fiber.Ctx) error {
	menuID := c.Params("id")

	var req dto.UpdateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	grpcRequest := &product.UpdateMenuRequest{
		MenuId:      menuID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageUrl:    req.ImageURL,
		CategoryId:  req.CategoryID,
	}

	res, err := h.productClient.UpdateMenu(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// DeleteMenu menangani permintaan HTTP untuk menghapus menu.
func (h *productHandler) DeleteMenu(c *fiber.Ctx) error {
	menuID := c.Params("id")

	grpcRequest := &product.DeleteMenuRequest{MenuId: menuID}

	res, err := h.productClient.DeleteMenu(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// ============== CATEGORY HANDLERS ==============

// CreateCategory menangani permintaan HTTP untuk membuat kategori baru.
func (h *productHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	// Validasi bahwa StoreID ada di body dan merupakan UUID yang valid.
	if _, err := uuid.Parse(req.StoreID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or missing storeId in request body",
		})
	}

	grpcRequest := &product.CreateCategoryRequest{
		Name:    req.Name,
		StoreId: req.StoreID,
	}

	res, err := h.productClient.CreateCategory(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetCategoryByID menangani permintaan HTTP untuk mencari satu kategori berdasarkan ID.
func (h *productHandler) GetCategoryByID(c *fiber.Ctx) error {
	categoryID := c.Params("id")

	grpcRequest := &product.GetCategoryByIDRequest{CategoryId: categoryID}

	res, err := h.productClient.GetCategoryByID(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateCategory menangani permintaan HTTP untuk memperbarui kategori.
func (h *productHandler) UpdateCategory(c *fiber.Ctx) error {
	categoryID := c.Params("id")

	var req dto.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	grpcRequest := &product.UpdateCategoryRequest{
		CategoryId: categoryID,
		Name:       req.Name,
	}

	res, err := h.productClient.UpdateCategory(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// DeleteCategory menangani permintaan HTTP untuk menghapus kategori.
func (h *productHandler) DeleteCategory(c *fiber.Ctx) error {
	categoryID := c.Params("id")

	grpcRequest := &product.DeleteCategoryRequest{CategoryId: categoryID}

	res, err := h.productClient.DeleteCategory(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// ============== POS HANDLERS ==============

// GetMenusByStoreCode mengambil menu berdasarkan store code
func (h *productHandler) GetMenusByStoreCode(c *fiber.Ctx) error {
	storeCode := c.Params("storeCode")
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	// 1. Panggil store-service untuk mendapatkan ID toko
	storeRes, err := h.storeClient.GetStoreByCode(ctx, &store.GetStoreByCodeRequest{StoreCode: storeCode})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "store not found"})
	}

	// 2. Gunakan ID toko untuk mengambil menu dari product-service
	menusRes, err := h.productClient.GetMenusByStoreID(ctx, &product.GetMenusByStoreIDRequest{StoreId: storeRes.Store.Id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get menus for the store"})
	}

	return c.Status(fiber.StatusOK).JSON(menusRes)
}

// GetMenusByStoreID mengambil menu berdasarkan ID toko
func (h *productHandler) GetMenusByStoreID(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	menusRes, err := h.productClient.GetMenusByStoreID(ctx, &product.GetMenusByStoreIDRequest{StoreId: id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get menus"})
	}

	return c.Status(fiber.StatusOK).JSON(menusRes)
}

// GetCategoriesByStoreCode mengambil kategori berdasarkan store code
func (h *productHandler) GetCategoriesByStoreCode(c *fiber.Ctx) error {
	storeCode := c.Params("storeCode")
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	// 1. Panggil store-service untuk mendapatkan ID toko
	storeRes, err := h.storeClient.GetStoreByCode(ctx, &store.GetStoreByCodeRequest{StoreCode: storeCode})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "store not found"})
	}

	// 2. Gunakan ID toko untuk mengambil kategori dari product-service
	categoriesRes, err := h.productClient.GetCategoriesByStoreID(ctx, &product.GetCategoriesByStoreIDRequest{StoreId: storeRes.Store.Id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get categories for the store"})
	}

	return c.Status(fiber.StatusOK).JSON(categoriesRes)
}

// GetCategoriesByStoreID mengambil kategori berdasarkan ID toko
func (h *productHandler) GetCategoriesByStoreID(c *fiber.Ctx) error {
	id := c.Params("id")

	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	defer cancel()

	categoriesRes, err := h.productClient.GetCategoriesByStoreID(ctx, &product.GetCategoriesByStoreIDRequest{StoreId: id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get categories"})
	}
	return c.Status(fiber.StatusOK).JSON(categoriesRes)
}
