package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	// Blok import yang diperbaiki dan lengkap
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/dto"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	orderPb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
)

type OrderHandler struct {
	client client.IOrderServiceClient
}

func NewOrderHandler(client client.IOrderServiceClient) *OrderHandler {
	return &OrderHandler{client: client}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	body := new(dto.CreateOrderRequest)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse user ID"})
	}

	grpcRequest := &orderPb.CreateOrderRequest{
		UserId:  userIDInt,
		StoreId: int64(body.StoreID),
	}

	for _, item := range body.Items {
		grpcRequest.Items = append(grpcRequest.Items, &orderPb.OrderItemRequest{
			ProductId: int64(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}

	res, err := h.client.GetOrderServiceClient().CreateOrder(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	orderID, err := strconv.ParseInt(c.Params("orderID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	grpcRequest := &orderPb.GetOrderRequest{Id: orderID}

	res, err := h.client.GetOrderServiceClient().GetOrder(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in token"})
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse user ID"})
	}

	grpcRequest := &orderPb.GetAllOrdersRequest{UserId: userIDInt}

	res, err := h.client.GetOrderServiceClient().GetAllOrders(c.Context(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
