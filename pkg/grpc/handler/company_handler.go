package handler

import (
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	// companyPb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
)

type CompanyHandler struct {
	client client.ICompanyServiceClient // Gunakan interface
}

func NewCompanyHandler(client client.ICompanyServiceClient) *CompanyHandler { // Terima interface
	return &CompanyHandler{client: client}
}
