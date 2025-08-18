// internal/product/handler.go
package product

import (
	"fnb-system/pkg/dto"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{service}
}

// --- Category Handlers ---

func (h *ProductHandler) CreateCategory(c *fiber.Ctx) error {
	var req dto.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Anda dapat menambahkan validasi untuk request body di sini

	category, err := h.service.CreateCategory(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

func (h *ProductHandler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(categories)
}

// --- Menu Handlers ---

func (h *ProductHandler) CreateMenu(c *fiber.Ctx) error {
	var req dto.MenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Anda dapat menambahkan validasi untuk request body di sini

	menu, err := h.service.CreateMenu(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(menu)
}

func (h *ProductHandler) GetAllMenus(c *fiber.Ctx) error {
	menus, err := h.service.GetAllMenus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(menus)
}

// GetMenuByID menangani permintaan untuk mendapatkan satu menu berdasarkan ID.
// Ini adalah fungsi yang kita tambahkan untuk memperbaiki error sebelumnya.
func (h *ProductHandler) GetMenuByID(c *fiber.Ctx) error {
	// Ambil ID dari parameter URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	// Panggil service untuk mendapatkan data menu
	menu, err := h.service.GetMenuByID(uint(id))
	if err != nil {
		// Kemungkinan besar menu tidak ditemukan
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Menu not found"})
	}

	return c.Status(fiber.StatusOK).JSON(menu)
}
