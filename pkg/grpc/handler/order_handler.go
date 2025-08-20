package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/dto"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
)

// OrderHandler mendefinisikan interface untuk semua metode handler pesanan.
type OrderHandler interface {
	CreateOrder(c *fiber.Ctx) error
	GetOrder(c *fiber.Ctx) error
	CalculateBill(c *fiber.Ctx) error
	CreatePublicOrder(c *fiber.Ctx) error
	UpdateOrderStatus(c *fiber.Ctx) error
	UpdateOrderItems(c *fiber.Ctx) error
	ConfirmPayment(c *fiber.Ctx) error
}

type orderHandler struct {
	client pb.OrderServiceClient
}

// NewOrderHandler adalah konstruktor untuk orderHandler.
func NewOrderHandler(client pb.OrderServiceClient) OrderHandler {
	return &orderHandler{client: client}
}

// CreateOrder menangani pembuatan pesanan oleh pengguna yang terautentikasi.
func (h *orderHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found in token"})
	}

	var grpcOrderItems []*pb.OrderItemRequest
	for _, item := range req.Items {
		grpcOrderItems = append(grpcOrderItems, &pb.OrderItemRequest{
			ProductId: strconv.FormatUint(uint64(item.ProductID), 10),
			Quantity:  int32(item.Quantity),
		})
	}

	grpcRequest := &pb.CreateOrderRequest{
		UserId:  userIDStr,
		StoreId: strconv.FormatUint(uint64(req.StoreID), 10),
		Items:   grpcOrderItems,
	}

	res, err := h.client.CreateOrder(context.Background(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetOrder mengambil detail pesanan berdasarkan ID.
func (h *orderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	req := &pb.GetOrderRequest{Id: id}

	res, err := h.client.GetOrder(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// CalculateBill menghitung total tagihan dari item-item yang diberikan.
func (h *orderHandler) CalculateBill(c *fiber.Ctx) error {
	var req pb.CalculateBillRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.client.CalculateBill(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// CreatePublicOrder menangani pembuatan pesanan oleh publik (tanpa login).
func (h *orderHandler) CreatePublicOrder(c *fiber.Ctx) error {
	storeCode := c.Params("storeCode")
	var req pb.CreatePublicOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req.StoreCode = storeCode

	res, err := h.client.CreatePublicOrder(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

// UpdateOrderStatus memperbarui status sebuah pesanan.
func (h *orderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req pb.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req.Id = id

	res, err := h.client.UpdateOrderStatus(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// UpdateOrderItems mengubah item-item di dalam pesanan yang ada.
func (h *orderHandler) UpdateOrderItems(c *fiber.Ctx) error {
	id := c.Params("id")
	var req pb.UpdateOrderItemsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req.Id = id

	res, err := h.client.UpdateOrderItems(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// ConfirmPayment mengonfirmasi pembayaran untuk sebuah pesanan.
func (h *orderHandler) ConfirmPayment(c *fiber.Ctx) error {
	var req pb.ConfirmPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := h.client.ConfirmPayment(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
