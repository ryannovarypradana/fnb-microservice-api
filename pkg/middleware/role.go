package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

// Authorize is a middleware to check user roles against a list of allowed roles.
func Authorize(allowedRoles ...model.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get the role from c.Locals, which was set by the AuthMiddleware.
		role, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "Authorization information not found in token.",
			})
		}

		// 2. Check if the user's role is in the list of allowed roles.
		isAllowed := false
		for _, allowedRole := range allowedRoles {
			if model.Role(role) == allowedRole {
				isAllowed = true
				break
			}
		}

		// 3. If the role is not allowed, return a 403 Forbidden error.
		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": true,
				"msg":   "You do not have permission to access this resource.",
			})
		}

		// 4. If the role is allowed, proceed to the next handler.
		return c.Next()
	}
}
