package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
)

type StoreHandler struct {
	client store.StoreServiceClient
}

func NewStoreHandler(client store.StoreServiceClient) *StoreHandler {
	return &StoreHandler{client: client}
}

func (h *StoreHandler) CreateStore(c *fiber.Ctx) error {
	req := new(store.CreateStoreRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Di aplikasi nyata, CompanyId bisa diambil dari token JWT (actor)
	// atau divalidasi berdasarkan hak akses pengguna.
	// Misalnya: claims := c.Locals("user_claims").(jwt.MapClaims)
	// req.CompanyId = claims["company_id"].(string)

	res, err := h.client.CreateStore(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create store"})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *StoreHandler) GetStore(c *fiber.Ctx) error {
	req := &store.GetStoreRequest{
		Id: c.Params("id"),
	}

	res, err := h.client.GetStore(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "store not found"})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *StoreHandler) GetAllStores(c *fiber.Ctx) error {
	req := &store.GetAllStoresRequest{
		Search: c.Query("search"),
	}

	res, err := h.client.GetAllStores(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
