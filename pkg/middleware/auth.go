package middleware

import (
	"fnb-system/internal/auth" // Kita butuh fungsi ValidateToken dari service auth
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware adalah middleware untuk memverifikasi token JWT.
func AuthMiddleware(c *fiber.Ctx) error {
	// 1. Ambil header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	// 2. Cek format header "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	tokenString := parts[1]

	// 3. Validasi token
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// 4. Simpan informasi dari token ke dalam context request
	// agar bisa digunakan oleh handler selanjutnya.
	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)

	// 5. Lanjutkan ke handler/middleware berikutnya
	return c.Next()
}
