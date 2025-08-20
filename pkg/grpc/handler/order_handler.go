// pkg/grpc/handler/order_handler.go

package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/dto"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
)

// OrderHandler adalah interface untuk handler pesanan.
type OrderHandler interface {
	CreateOrder(c *fiber.Ctx) error
}

type orderHandler struct {
	client order.OrderServiceClient
}

// NewOrderHandler membuat instance baru dari orderHandler.
func NewOrderHandler(client order.OrderServiceClient) OrderHandler {
	return &orderHandler{client: client}
}

// CreateOrder menangani permintaan HTTP untuk membuat pesanan baru.
func (h *orderHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in token",
		})
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid User ID in token",
		})
	}

	var grpcOrderItems []*order.OrderItemRequest
	for _, item := range req.Items {
		grpcOrderItems = append(grpcOrderItems, &order.OrderItemRequest{
			// PERBAIKAN FINAL: Menggunakan ProductId (dari proto) dan ProductID (dari DTO)
			ProductId: int64(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}

	grpcRequest := &order.CreateOrderRequest{
		UserId:  userID,
		StoreId: int64(req.StoreID),
		Items:   grpcOrderItems,
	}

	res, err := h.client.CreateOrder(context.Background(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}
