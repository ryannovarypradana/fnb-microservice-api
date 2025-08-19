// /pkg/middleware/auth.go
package middleware

import (
	"os"
	"strings"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model" // Pastikan pkg/model/role.go sudah dibuat

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware memeriksa validitas token JWT dari header Authorization.
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is required"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header format must be Bearer {token}"})
	}

	tokenString := parts[1]
	jwtSecret := os.Getenv("JWT_SECRET") // Pastikan .env menggunakan JWT_SECRET

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Menyimpan claims di context untuk digunakan oleh middleware selanjutnya atau handler
	c.Locals("user_claims", claims)

	return c.Next()
}

// Helper untuk mengambil claims yang sudah disimpan oleh AuthMiddleware
func getClaims(c *fiber.Ctx) (jwt.MapClaims, bool) {
	claims, ok := c.Locals("user_claims").(jwt.MapClaims)
	return claims, ok
}

// RoleMiddleware adalah middleware generik baru untuk memeriksa daftar role yang diizinkan.
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := getClaims(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token claims not found"})
		}

		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Invalid role format in token"})
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next() // Role diizinkan, lanjutkan
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: You do not have the required role"})
	}
}

// IsCompanyOwnerOrSuperAdmin adalah contoh middleware kompleks untuk otorisasi berbasis kepemilikan.
func IsCompanyOwnerOrSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := getClaims(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token claims not found"})
		}

		userRole := claims["role"].(string)

		// Super admin bisa melakukan apa saja
		if userRole == model.RoleSuperAdmin {
			return c.Next()
		}

		// Jika bukan super admin, harus company_rep
		if userRole == model.RoleCompanyRep {
			companyIDFromToken, ok := claims["company_id"].(string)
			if !ok {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Company ID not found in token for company representative"})
			}

			companyIDFromParam := c.Params("companyId") // Contoh: /api/companies/:companyId/stores

			if companyIDFromToken == companyIDFromParam {
				return c.Next() // Token milik company yang sama
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: You do not own this resource"})
	}
}
