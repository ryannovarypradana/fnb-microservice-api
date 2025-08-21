package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/dto"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
)

// CompanyHandler adalah interface untuk handler perusahaan.
type CompanyHandler interface {
	CreateCompany(c *fiber.Ctx) error
	GetCompany(c *fiber.Ctx) error
	GetAllCompanies(c *fiber.Ctx) error // <-- DITAMBAHKAN
}

type companyHandler struct {
	client company.CompanyServiceClient
}

// NewCompanyHandler membuat instance baru dari companyHandler.
func NewCompanyHandler(client company.CompanyServiceClient) CompanyHandler {
	return &companyHandler{client: client}
}

// CreateCompany menangani permintaan HTTP untuk membuat perusahaan baru.
func (h *companyHandler) CreateCompany(c *fiber.Ctx) error {
	var req dto.CreateCompanyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request body",
		})
	}

	grpcRequest := &company.CreateCompanyRequest{
		Name:    req.Name,
		Address: req.Address,
	}

	res, err := h.client.CreateCompany(context.Background(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// GetCompany menangani permintaan HTTP untuk mendapatkan detail perusahaan berdasarkan ID.
func (h *companyHandler) GetCompany(c *fiber.Ctx) error {
	companyID := c.Params("id")
	if companyID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Company ID is required",
		})
	}

	grpcRequest := &company.GetCompanyRequest{
		Id: companyID,
	}

	res, err := h.client.GetCompany(context.Background(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetAllCompanies menangani permintaan untuk mendapatkan semua perusahaan.
func (h *companyHandler) GetAllCompanies(c *fiber.Ctx) error {
	searchQuery := c.Query("search") // Mengambil query param "search"

	grpcRequest := &company.GetAllCompaniesRequest{
		Search: searchQuery,
	}

	res, err := h.client.GetAllCompanies(context.Background(), grpcRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
