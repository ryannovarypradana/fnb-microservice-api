package handler

import (
	"github.com/gofiber/fiber/v2"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StoreHandler struct {
	client pb.StoreServiceClient
}

func NewStoreHandler(client pb.StoreServiceClient) *StoreHandler {
	return &StoreHandler{client: client}
}

func (h *StoreHandler) CreateStore(c *fiber.Ctx) error {
	var req pb.CreateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	companyID, ok := c.Locals("company_id").(string)
	if !ok || companyID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid company id"})
	}

	req.CompanyId = companyID

	res, err := h.client.CreateStore(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (h *StoreHandler) GetStore(c *fiber.Ctx) error {
	id := c.Params("id")
	req := &pb.GetStoreRequest{Id: id}

	res, err := h.client.GetStore(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// --- FUNGSI BARU DIMULAI DI SINI ---
func (h *StoreHandler) GetStoreByCode(c *fiber.Ctx) error {
	storeCode := c.Params("storeCode")
	req := &pb.GetStoreByCodeRequest{StoreCode: storeCode}

	res, err := h.client.GetStoreByCode(c.Context(), req)
	if err != nil {
		// Handle error specifically for not found cases from gRPC status
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": st.Message()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if the store object is actually present in the response
	if res == nil || res.Store == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "store not found"})
	}

	return c.JSON(res.Store)
}

// --- FUNGSI BARU SELESAI DI SINI ---

func (h *StoreHandler) GetAllStores(c *fiber.Ctx) error {
	searchQuery := c.Query("search")
	req := &pb.GetAllStoresRequest{Search: searchQuery}

	res, err := h.client.GetAllStores(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res.Stores)
}

func (h *StoreHandler) UpdateStore(c *fiber.Ctx) error {
	id := c.Params("id")
	var req pb.UpdateStoreRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	req.Id = id

	res, err := h.client.UpdateStore(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *StoreHandler) DeleteStore(c *fiber.Ctx) error {
	id := c.Params("id")
	req := &pb.DeleteStoreRequest{Id: id}

	res, err := h.client.DeleteStore(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *StoreHandler) CloneStoreContent(c *fiber.Ctx) error {
	var req pb.CloneStoreContentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.SourceStoreId == "" || req.DestinationStoreId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "source_store_id and destination_store_id are required"})
	}

	res, err := h.client.CloneStoreContent(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}
