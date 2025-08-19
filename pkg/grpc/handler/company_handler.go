package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
)

type CompanyHandler struct {
	client company.CompanyServiceClient
}

func NewCompanyHandler(client company.CompanyServiceClient) *CompanyHandler {
	return &CompanyHandler{client: client}
}

func (h *CompanyHandler) GetAllCompanies(c *fiber.Ctx) error {
	req := &company.GetAllCompaniesRequest{
		Search: c.Query("search"),
	}
	res, err := h.client.GetAllCompanies(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *CompanyHandler) GetCompany(c *fiber.Ctx) error {
	req := &company.GetCompanyRequest{
		Id: c.Params("companyId"),
	}
	res, err := h.client.GetCompany(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "company not found"})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
