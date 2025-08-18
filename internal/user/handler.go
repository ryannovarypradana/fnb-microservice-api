// internal/user/handler.go
package user

import (
	"fnb-system/pkg/dto"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService}
}

// GetAllUsers menangani request untuk mendapatkan semua pengguna.
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Setup paginasi dari query params
	pagination := dto.Pagination{
		Limit: 10,
		Page:  1,
		Sort:  "id asc",
	} // Anda bisa menambahkan logika untuk mengambil nilai ini dari c.Query()

	users, err := h.userService.GetAll(&pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

// GetUserByID menangani request untuk mendapatkan satu pengguna.
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	// Ambil ID dari parameter URL dan konversi ke integer
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdateUser menangani request untuk memperbarui pengguna.
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	var req dto.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Anda bisa menambahkan validasi DTO di sini

	user, err := h.userService.Update(uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// DeleteUser menangani request untuk menghapus pengguna.
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	err = h.userService.Delete(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent) // 204 No Content untuk delete sukses
}
