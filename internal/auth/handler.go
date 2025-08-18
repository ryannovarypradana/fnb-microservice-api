package auth

import (
	"fnb-system/pkg/dto"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

// Register menangani request registrasi user.
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.AuthRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Anda bisa menambahkan validasi request di sini

	user, err := h.authService.Register(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Buat response DTO agar tidak mengirim password hash
	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// Login menangani request login.
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.AuthLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Anda bisa menambahkan validasi request di sini

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}
