// /pkg/grpc/handler/user_handler.go
package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
)

// UserHandler connects the API Gateway to the User gRPC service.
type UserHandler struct {
	client user.UserServiceClient
}

// NewUserHandler creates a new handler for user routes.
func NewUserHandler(client user.UserServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

// GetUser handles the GET /users/:id route.
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	req := &user.GetUserRequest{
		Id: c.Params("id"),
	}

	res, err := h.client.GetUser(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
