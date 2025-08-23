package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
)

type AuthHandler struct {
	client pb.AuthServiceClient
}

func NewAuthHandler(client pb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{client: client}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {

	log.Println("====== LOGIN ATTEMPT RECEIVED ======")
	log.Printf("Request Body: %s", string(c.Body()))
	log.Printf("Content-Type Header: %s", c.Get("Content-Type"))
	log.Println("====================================")
	// ==================================================================
	req := new(pb.LoginRequest)
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
	req := new(pb.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.client.Register(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// RegisterStaff handles the creation of new staff members by authorized users.
func (h *AuthHandler) RegisterStaff(c *fiber.Ctx) error {
	req := new(pb.RegisterStaffRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// The gRPC call to the auth-service is made here.
	// The auth-service contains the core logic for hashing the password
	// and creating the user in the database with the specified role.
	res, err := h.client.RegisterStaff(c.Context(), req)
	if err != nil {
		// Errors from the gRPC service are propagated here.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}
