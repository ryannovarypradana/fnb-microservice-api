package store

import (
	"context"

	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StoreGRPCHandler struct {
	pb.UnimplementedStoreServiceServer
	service StoreService
}

func NewStoreGRPCHandler(grpcServer *grpc.Server, service StoreService) {
	handler := &StoreGRPCHandler{service: service}
	pb.RegisterStoreServiceServer(grpcServer, handler)
}

func (h *StoreGRPCHandler) CreateStore(ctx context.Context, req *pb.CreateStoreRequest) (*pb.CreateStoreResponse, error) {
	store, err := h.service.CreateStore(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create store: %v", err)
	}
	return &pb.CreateStoreResponse{Store: &pb.Store{
		Id:        store.ID.String(),
		CompanyId: store.CompanyID.String(),
		Name:      store.Name,
		Address:   store.Location,
		Code:      store.Code,
	}}, nil
}

func (h *StoreGRPCHandler) GetStore(ctx context.Context, req *pb.GetStoreRequest) (*pb.GetStoreResponse, error) {
	store, err := h.service.GetStore(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "store not found: %v", err)
	}
	return &pb.GetStoreResponse{Store: &pb.Store{
		Id:        store.ID.String(),
		CompanyId: store.CompanyID.String(),
		Name:      store.Name,
		Address:   store.Location,
		Code:      store.Code,
	}}, nil
}

func (h *StoreGRPCHandler) GetAllStores(ctx context.Context, req *pb.GetAllStoresRequest) (*pb.GetAllStoresResponse, error) {
	stores, err := h.service.GetAllStores(ctx, req.Search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get stores: %v", err)
	}

	var pbStores []*pb.Store
	for _, store := range stores {
		pbStores = append(pbStores, &pb.Store{
			Id:        store.ID.String(),
			CompanyId: store.CompanyID.String(),
			Name:      store.Name,
			Address:   store.Location,
			Code:      store.Code,
		})
	}

	return &pb.GetAllStoresResponse{Stores: pbStores}, nil
}

func (h *StoreGRPCHandler) UpdateStore(ctx context.Context, req *pb.UpdateStoreRequest) (*pb.GetStoreResponse, error) {
	store, err := h.service.UpdateStore(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update store: %v", err)
	}
	return &pb.GetStoreResponse{Store: &pb.Store{
		Id:        store.ID.String(),
		CompanyId: store.CompanyID.String(),
		Name:      store.Name,
		Address:   store.Location,
		Code:      store.Code,
	}}, nil
}

func (h *StoreGRPCHandler) DeleteStore(ctx context.Context, req *pb.DeleteStoreRequest) (*pb.DeleteStoreResponse, error) {
	if err := h.service.DeleteStore(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete store: %v", err)
	}
	return &pb.DeleteStoreResponse{Message: "Store deleted successfully"}, nil
}

func (h *StoreGRPCHandler) CloneStoreContent(ctx context.Context, req *pb.CloneStoreContentRequest) (*pb.CloneStoreContentResponse, error) {
	if err := h.service.CloneStoreContent(ctx, req.SourceStoreId, req.DestinationStoreId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clone store content: %v", err)
	}
	return &pb.CloneStoreContentResponse{Message: "Store content cloned successfully"}, nil
}
