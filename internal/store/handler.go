package store

import (
	"context"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
)

type GRPCHandler struct {
	store.UnimplementedStoreServiceServer
	service Service
}

// NewGRPCHandler creates a new gRPC handler for the store service.
func NewGRPCHandler(s Service) *GRPCHandler {
	return &GRPCHandler{service: s}
}

func (h *GRPCHandler) CreateStore(ctx context.Context, req *store.CreateStoreRequest) (*store.CreateStoreResponse, error) {
	createdStore, err := h.service.CreateStore(ctx, req.GetCompanyId(), req.GetName(), req.GetAddress())
	if err != nil {
		return nil, err
	}

	return &store.CreateStoreResponse{
		Store: &store.Store{
			Id:        createdStore.ID.String(),
			CompanyId: createdStore.CompanyID.String(),
			Name:      createdStore.Name,
			Address:   createdStore.Address,
			Code:      createdStore.Code,
		},
	}, nil
}

func (h *GRPCHandler) GetStore(ctx context.Context, req *store.GetStoreRequest) (*store.GetStoreResponse, error) {
	foundStore, err := h.service.GetStoreByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &store.GetStoreResponse{
		Store: &store.Store{
			Id:        foundStore.ID.String(),
			CompanyId: foundStore.CompanyID.String(),
			Name:      foundStore.Name,
			Address:   foundStore.Address,
			Code:      foundStore.Code,
		},
	}, nil
}

func (h *GRPCHandler) GetAllStores(ctx context.Context, req *store.GetAllStoresRequest) (*store.GetAllStoresResponse, error) {
	stores, err := h.service.GetAllStores(ctx, req.GetSearch())
	if err != nil {
		return nil, err
	}

	var storeMessages []*store.Store
	for _, s := range stores {
		storeMessages = append(storeMessages, &store.Store{
			Id:        s.ID.String(),
			CompanyId: s.CompanyID.String(),
			Name:      s.Name,
			Address:   s.Address,
			Code:      s.Code,
		})
	}

	return &store.GetAllStoresResponse{Stores: storeMessages}, nil
}
