package company

import (
	"context"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
)

type GRPCHandler struct {
	company.UnimplementedCompanyServiceServer
	service Service
}

func NewGRPCHandler(s Service) *GRPCHandler {
	return &GRPCHandler{service: s}
}

func (h *GRPCHandler) CreateCompany(ctx context.Context, req *company.CreateCompanyRequest) (*company.CreateCompanyResponse, error) {
	c, err := h.service.CreateCompany(ctx, req.GetName(), req.GetAddress())
	if err != nil {
		return nil, err
	}

	return &company.CreateCompanyResponse{
		Company: &company.Company{
			Id:      c.ID.String(),
			Name:    c.Name,
			Address: c.Address,
			Code:    c.Code,
		},
	}, nil
}

func (h *GRPCHandler) GetCompany(ctx context.Context, req *company.GetCompanyRequest) (*company.GetCompanyResponse, error) {
	c, err := h.service.GetCompanyByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &company.GetCompanyResponse{
		Company: &company.Company{
			Id:      c.ID.String(),
			Name:    c.Name,
			Address: c.Address,
			Code:    c.Code,
		},
	}, nil
}

func (h *GRPCHandler) GetAllCompanies(ctx context.Context, req *company.GetAllCompaniesRequest) (*company.GetAllCompaniesResponse, error) {
	companies, err := h.service.GetAllCompanies(ctx, req.GetSearch())
	if err != nil {
		return nil, err
	}
	var companyMessages []*company.Company
	for _, c := range companies {
		companyMessages = append(companyMessages, &company.Company{
			Id:      c.ID.String(),
			Name:    c.Name,
			Address: c.Address,
			Code:    c.Code,
		})
	}
	return &company.GetAllCompaniesResponse{Companies: companyMessages}, nil
}
