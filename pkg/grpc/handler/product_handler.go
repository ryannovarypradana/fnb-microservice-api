package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
)

type ProductHandler struct {
	client product.ProductServiceClient
}

func NewProductHandler(client product.ProductServiceClient) *ProductHandler {
	return &ProductHandler{client: client}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	req := new(product.CreateProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Anda bisa menambahkan logika untuk mengambil StoreID dari token JWT di sini jika diperlukan
	// claims := c.Locals("user_claims").(jwt.MapClaims)
	// req.StoreId = claims["store_id"].(string)

	res, err := h.client.CreateProduct(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create product"})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	req := &product.GetProductRequest{
		Id: c.Params("id"),
	}

	res, err := h.client.GetProduct(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	req := &product.GetAllProductsRequest{
		Search: c.Query("search"),
		// Anda bisa menambahkan filter lain dari query params, misal: "store_id"
		// StoreId: c.Query("store_id"),
	}

	res, err := h.client.GetAllProducts(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to retrieve products"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	req := new(product.UpdateProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Set ID dari URL parameter
	req.Id = c.Params("id")

	res, err := h.client.UpdateProduct(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update product"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	req := &product.DeleteProductRequest{
		Id: c.Params("id"),
	}

	res, err := h.client.DeleteProduct(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete product"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
