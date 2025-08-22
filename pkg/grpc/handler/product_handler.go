// pkg/grpc/handler/product_handler.go

package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/dto"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
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
}

type productHandler struct {
	client product.ProductServiceClient
}

// NewProductHandler membuat instance baru dari productHandler.
func NewProductHandler(client product.ProductServiceClient) ProductHandler {
	return &productHandler{client: client}
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

	grpcRequest := &product.CreateMenuRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageUrl:    req.ImageURL,
		CategoryId:  int64(req.CategoryID),
		StoreId:     int64(req.StoreID),
	}

	res, err := h.client.CreateMenu(context.Background(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetMenuByID menangani permintaan HTTP untuk mencari satu menu berdasarkan ID.
func (h *productHandler) GetMenuByID(c *fiber.Ctx) error {
	menuID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Menu ID",
		})
	}

	grpcRequest := &product.GetMenuByIDRequest{
		MenuId: menuID,
	}

	res, err := h.client.GetMenuByID(context.Background(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateMenu menangani permintaan HTTP untuk memperbarui menu.
func (h *productHandler) UpdateMenu(c *fiber.Ctx) error {
	menuID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Menu ID"})
	}

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
		CategoryId:  int64(req.CategoryID),
	}

	res, err := h.client.UpdateMenu(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// DeleteMenu menangani permintaan HTTP untuk menghapus menu.
func (h *productHandler) DeleteMenu(c *fiber.Ctx) error {
	menuID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Menu ID"})
	}

	grpcRequest := &product.DeleteMenuRequest{MenuId: menuID}

	res, err := h.client.DeleteMenu(c.Context(), grpcRequest)
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

	grpcRequest := &product.CreateCategoryRequest{
		Name:    req.Name,
		StoreId: int64(req.StoreID),
	}

	res, err := h.client.CreateCategory(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetCategoryByID menangani permintaan HTTP untuk mencari satu kategori berdasarkan ID.
func (h *productHandler) GetCategoryByID(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Category ID"})
	}

	grpcRequest := &product.GetCategoryByIDRequest{CategoryId: categoryID}

	res, err := h.client.GetCategoryByID(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdateCategory menangani permintaan HTTP untuk memperbarui kategori.
func (h *productHandler) UpdateCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Category ID"})
	}

	var req dto.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request body"})
	}

	grpcRequest := &product.UpdateCategoryRequest{
		CategoryId: categoryID,
		Name:       req.Name,
	}

	res, err := h.client.UpdateCategory(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// DeleteCategory menangani permintaan HTTP untuk menghapus kategori.
func (h *productHandler) DeleteCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Category ID"})
	}

	grpcRequest := &product.DeleteCategoryRequest{CategoryId: categoryID}

	res, err := h.client.DeleteCategory(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
