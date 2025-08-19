package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
)

type OrderHandler struct {
	client order.OrderServiceClient
}

func NewOrderHandler(client order.OrderServiceClient) *OrderHandler {
	return &OrderHandler{client: client}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	req := new(order.CreateOrderRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	claims, ok := c.Locals("user_claims").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	userId, ok := claims["id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user id not found in token"})
	}

	req.UserId = userId

	res, err := h.client.CreateOrder(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *OrderHandler) GetOrderById(c *fiber.Ctx) error {
	req := &order.GetOrderByIdRequest{
		Id: c.Params("id"),
	}

	// Logika otorisasi tambahan bisa dilakukan di sini atau di service-level
	// Misalnya, memastikan user yang request adalah pemilik order atau seorang admin.

	res, err := h.client.GetOrderById(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "order not found"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *OrderHandler) GetAllOrdersForUser(c *fiber.Ctx) error {
	claims, ok := c.Locals("user_claims").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	userId, ok := claims["id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user id not found in token"})
	}

	req := &order.GetAllOrdersRequest{
		UserId: userId, // Filter pesanan hanya untuk pengguna yang sedang login
	}

	res, err := h.client.GetAllOrders(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
