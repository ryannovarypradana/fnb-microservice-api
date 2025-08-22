package handler

import (
	"github.com/gofiber/fiber/v2"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/promotion"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PromotionHandler mendefinisikan interface untuk semua metode handler promosi.
type PromotionHandler interface {
	// Discount Handlers
	CreateDiscount(c *fiber.Ctx) error
	GetDiscount(c *fiber.Ctx) error
	ListDiscounts(c *fiber.Ctx) error
	UpdateDiscount(c *fiber.Ctx) error
	DeleteDiscount(c *fiber.Ctx) error

	// Voucher Handlers
	CreateVoucher(c *fiber.Ctx) error
	GetVoucher(c *fiber.Ctx) error
	ListVouchers(c *fiber.Ctx) error
	UpdateVoucher(c *fiber.Ctx) error
	DeleteVoucher(c *fiber.Ctx) error

	// Bundle Handlers
	CreateBundle(c *fiber.Ctx) error
	GetBundle(c *fiber.Ctx) error
	ListBundles(c *fiber.Ctx) error
	UpdateBundle(c *fiber.Ctx) error
	DeleteBundle(c *fiber.Ctx) error
}

type promotionHandler struct {
	client pb.PromotionServiceClient
}

// NewPromotionHandler adalah konstruktor untuk promotionHandler.
func NewPromotionHandler(client pb.PromotionServiceClient) PromotionHandler {
	return &promotionHandler{client: client}
}

// ========== Discount Handlers ==========

func (h *promotionHandler) CreateDiscount(c *fiber.Ctx) error {
	var req pb.CreateDiscountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.client.CreateDiscount(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *promotionHandler) GetDiscount(c *fiber.Ctx) error {
	// Diperbaiki: Menggunakan ID sebagai string, bukan uint64
	id := c.Params("id")
	res, err := h.client.GetDiscount(c.Context(), &pb.GetByIDRequest{Id: id})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *promotionHandler) ListDiscounts(c *fiber.Ctx) error {
	res, err := h.client.ListDiscounts(c.Context(), &emptypb.Empty{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *promotionHandler) UpdateDiscount(c *fiber.Ctx) error {
	// Diperbaiki: Menggunakan ID sebagai string, bukan uint64
	id := c.Params("id")
	var req pb.UpdateDiscountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req.Id = id
	res, err := h.client.UpdateDiscount(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *promotionHandler) DeleteDiscount(c *fiber.Ctx) error {
	// Diperbaiki: Menggunakan ID sebagai string, bukan uint64
	id := c.Params("id")
	res, err := h.client.DeleteDiscount(c.Context(), &pb.GetByIDRequest{Id: id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// ========== Voucher Handlers ==========

func (h *promotionHandler) CreateVoucher(c *fiber.Ctx) error {
	var req pb.CreateVoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.client.CreateVoucher(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *promotionHandler) GetVoucher(c *fiber.Ctx) error {
	// Diperbaiki: Menggunakan ID sebagai string, bukan uint64
	id := c.Params("id")
	res, err := h.client.GetVoucher(c.Context(), &pb.GetByIDRequest{Id: id})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *promotionHandler) ListVouchers(c *fiber.Ctx) error {
	res, err := h.client.ListVouchers(c.Context(), &emptypb.Empty{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *promotionHandler) UpdateVoucher(c *fiber.Ctx) error {
	// Diperbaiki: Menggunakan ID sebagai string, bukan uint64
	id := c.Params("id")
	var req pb.UpdateVoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req.Id = id
	res, err := h.client.UpdateVoucher(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *promotionHandler) DeleteVoucher(c *fiber.Ctx) error {
	// Diperbaiki: Menggunakan ID sebagai string, bukan uint64
	id := c.Params("id")
	res, err := h.client.DeleteVoucher(c.Context(), &pb.GetByIDRequest{Id: id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// ========== Bundle Handlers ==========

func (h *promotionHandler) CreateBundle(c *fiber.Ctx) error {
	var req pb.CreateBundleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	res, err := h.client.CreateBundle(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *promotionHandler) GetBundle(c *fiber.Ctx) error {
	// Diperbaiki: Menggunakan ID sebagai string, bukan uint64
	id := c.Params("id")
	res, err := h.client.GetBundle(c.Context(), &pb.GetByIDRequest{Id: id})
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *promotionHandler) ListBundles(c *fiber.Ctx) error {
	res, err := h.client.ListBundles(c.Context(), &emptypb.Empty{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *promotionHandler) UpdateBundle(c *fiber.Ctx) error {
	// Diperbaiki: Menggunakan ID sebagai string, bukan uint64
	id := c.Params("id")
	var req pb.UpdateBundleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req.Id = id
	res, err := h.client.UpdateBundle(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *promotionHandler) DeleteBundle(c *fiber.Ctx) error {
	// Diperbaiki: Menggunakan ID sebagai string, bukan uint64
	id := c.Params("id")
	res, err := h.client.DeleteBundle(c.Context(), &pb.GetByIDRequest{Id: id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
