package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/dto"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
)

type ProductHandler struct {
	productClient client.IProductServiceClient
}

func NewProductHandler(productClient client.IProductServiceClient) *ProductHandler {
	return &ProductHandler{
		productClient: productClient,
	}
}

// --- Menu Handlers ---

// CreateMenu menangani permintaan untuk membuat menu baru.
func (h *ProductHandler) CreateMenu(c *fiber.Ctx) error {
	body := new(dto.CreateMenuRequest)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	grpcRequest := &pb.CreateMenuRequest{
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		ImageUrl:    body.ImageURL,
		CategoryId:  int64(body.CategoryID),
		StoreId:     int64(body.StoreID),
	}

	res, err := h.productClient.GetProductServiceClient().CreateMenu(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res.Menu)
}

// GetMenuByID mengambil detail menu berdasarkan ID.
func (h *ProductHandler) GetMenuByID(c *fiber.Ctx) error {
	menuID, err := strconv.ParseInt(c.Params("menuID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid menu ID format"})
	}

	grpcRequest := &pb.GetMenuByIDRequest{MenuId: menuID}
	res, err := h.productClient.GetProductServiceClient().GetMenuByID(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res.Menu)
}

// GetMenusByStoreID mengambil semua menu yang dimiliki oleh sebuah toko.
func (h *ProductHandler) GetMenusByStoreID(c *fiber.Ctx) error {
	storeID, err := strconv.ParseInt(c.Params("storeID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid store ID format"})
	}

	grpcRequest := &pb.GetMenusByStoreIDRequest{StoreId: storeID}
	res, err := h.productClient.GetProductServiceClient().GetMenusByStoreID(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res.Menus)
}

// UpdateMenu memperbarui data menu yang sudah ada.
func (h *ProductHandler) UpdateMenu(c *fiber.Ctx) error {
	menuID, err := strconv.ParseInt(c.Params("menuID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid menu ID format"})
	}

	body := new(dto.UpdateMenuRequest)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	grpcRequest := &pb.UpdateMenuRequest{
		MenuId:      menuID,
		Name:        body.Name,
		Description: body.Description,
		Price:       body.Price,
		ImageUrl:    body.ImageURL,
		CategoryId:  int64(body.CategoryID),
	}

	res, err := h.productClient.GetProductServiceClient().UpdateMenu(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res.Menu)
}

// DeleteMenu menghapus sebuah menu berdasarkan ID.
func (h *ProductHandler) DeleteMenu(c *fiber.Ctx) error {
	menuID, err := strconv.ParseInt(c.Params("menuID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid menu ID format"})
	}

	grpcRequest := &pb.DeleteMenuRequest{MenuId: menuID}
	res, err := h.productClient.GetProductServiceClient().DeleteMenu(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": res.Message})
}

// --- Category Handlers ---

// CreateCategory menangani permintaan untuk membuat kategori baru.
func (h *ProductHandler) CreateCategory(c *fiber.Ctx) error {
	body := new(dto.CreateCategoryRequest)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	grpcRequest := &pb.CreateCategoryRequest{
		Name:    body.Name,
		StoreId: int64(body.StoreID),
	}

	res, err := h.productClient.GetProductServiceClient().CreateCategory(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res.Category)
}

// GetCategoryByID mengambil detail kategori berdasarkan ID.
func (h *ProductHandler) GetCategoryByID(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseInt(c.Params("categoryID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID format"})
	}

	grpcRequest := &pb.GetCategoryByIDRequest{CategoryId: categoryID}
	res, err := h.productClient.GetProductServiceClient().GetCategoryByID(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res.Category)
}

// GetCategoriesByStoreID mengambil semua kategori yang dimiliki oleh sebuah toko.
func (h *ProductHandler) GetCategoriesByStoreID(c *fiber.Ctx) error {
	storeID, err := strconv.ParseInt(c.Params("storeID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid store ID format"})
	}

	grpcRequest := &pb.GetCategoriesByStoreIDRequest{StoreId: storeID}
	res, err := h.productClient.GetProductServiceClient().GetCategoriesByStoreID(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res.Categories)
}

// UpdateCategory memperbarui data kategori yang sudah ada.
func (h *ProductHandler) UpdateCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseInt(c.Params("categoryID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID format"})
	}

	body := new(dto.UpdateCategoryRequest)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body", "details": err.Error()})
	}

	grpcRequest := &pb.UpdateCategoryRequest{
		CategoryId: categoryID,
		Name:       body.Name,
	}

	res, err := h.productClient.GetProductServiceClient().UpdateCategory(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res.Category)
}

// DeleteCategory menghapus sebuah kategori berdasarkan ID.
func (h *ProductHandler) DeleteCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseInt(c.Params("categoryID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID format"})
	}

	grpcRequest := &pb.DeleteCategoryRequest{CategoryId: categoryID}
	res, err := h.productClient.GetProductServiceClient().DeleteCategory(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": res.Message})
}
