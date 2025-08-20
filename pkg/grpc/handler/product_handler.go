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
	CreateMenu(c *fiber.Ctx) error
	GetMenuByID(c *fiber.Ctx) error
	// Tambahkan method lain di sini nanti (CreateCategory, dll.)
}

type productHandler struct {
	client product.ProductServiceClient
}

// NewProductHandler membuat instance baru dari productHandler.
func NewProductHandler(client product.ProductServiceClient) ProductHandler {
	return &productHandler{client: client}
}

// CreateMenu menangani permintaan HTTP untuk membuat menu baru.
func (h *productHandler) CreateMenu(c *fiber.Ctx) error {
	// Menggunakan DTO yang benar dari file Anda: dto.CreateMenuRequest
	var req dto.CreateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	// Menggunakan message gRPC yang benar dari .proto: product.CreateMenuRequest
	grpcRequest := &product.CreateMenuRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageUrl:    req.ImageURL,
		CategoryId:  int64(req.CategoryID),
		StoreId:     int64(req.StoreID),
	}

	// Memanggil method gRPC yang benar: h.client.CreateMenu
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

	// Menggunakan message gRPC yang benar dari .proto: product.GetMenuByIDRequest
	grpcRequest := &product.GetMenuByIDRequest{
		MenuId: menuID,
	}

	// Memanggil method gRPC yang benar: h.client.GetMenuByID
	res, err := h.client.GetMenuByID(context.Background(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
