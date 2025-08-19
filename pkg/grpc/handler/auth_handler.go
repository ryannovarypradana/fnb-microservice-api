package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
)

type AuthHandler struct {
	client auth.AuthServiceClient
}

func NewAuthHandler(client auth.AuthServiceClient) *AuthHandler {
	return &AuthHandler{client: client}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(auth.LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.client.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(auth.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.client.Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}
