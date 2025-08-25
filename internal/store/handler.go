package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
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

func toProtoCompany(company *model.Company) *pb.Company {
	if company == nil || company.ID == uuid.Nil {
		return nil
	}
	return &pb.Company{
		Id:        company.ID.String(),
		Name:      company.Name,
		Code:      company.Code,
		CreatedAt: company.CreatedAt.Format(time.RFC3339),
		UpdatedAt: company.UpdatedAt.Format(time.RFC3339),
	}
}

// toProto adalah helper yang telah diperbarui untuk mengubah model.Store menjadi pb.Store dengan aman
func toProto(store *model.Store) *pb.Store {
	if store == nil {
		return nil
	}

	opHoursStr := ""
	if store.OperationalHours != nil {
		bytes, err := json.Marshal(store.OperationalHours)
		if err == nil {
			opHoursStr = string(bytes)
		}
	}

	// Pemeriksaan nilai nil untuk pointer
	var taxPercentage float32
	if store.TaxPercentage != nil {
		taxPercentage = float32(*store.TaxPercentage)
	}

	var latitude float64
	if store.Latitude != nil {
		latitude = *store.Latitude
	}

	var longitude float64
	if store.Longitude != nil {
		longitude = *store.Longitude
	}

	return &pb.Store{
		Id:               store.ID.String(),
		Name:             store.Name,
		Address:          store.Location,
		Code:             store.Code,
		BannerImageUrl:   store.BannerImageURL,
		TaxPercentage:    taxPercentage,
		OperationalHours: opHoursStr,
		Latitude:         latitude,
		Longitude:        longitude,
		Company:          toProtoCompany(&store.Company),
		CreatedAt:        store.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        store.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *StoreGRPCHandler) CreateStore(ctx context.Context, req *pb.CreateStoreRequest) (*pb.CreateStoreResponse, error) {
	store, err := h.service.CreateStore(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create store: %v", err)
	}
	return &pb.CreateStoreResponse{Store: toProto(store)}, nil
}

func (h *StoreGRPCHandler) GetStore(ctx context.Context, req *pb.GetStoreRequest) (*pb.GetStoreResponse, error) {
	store, err := h.service.GetStore(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "store not found: %v", err)
	}
	return &pb.GetStoreResponse{Store: toProto(store)}, nil
}

func (h *StoreGRPCHandler) GetAllStores(ctx context.Context, req *pb.GetAllStoresRequest) (*pb.GetAllStoresResponse, error) {
	stores, err := h.service.GetAllStores(ctx, req.Search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get stores: %v", err)
	}

	var pbStores []*pb.Store
	for _, store := range stores {
		pbStores = append(pbStores, toProto(store))
	}

	return &pb.GetAllStoresResponse{Stores: pbStores}, nil
}

func (h *StoreGRPCHandler) UpdateStore(ctx context.Context, req *pb.UpdateStoreRequest) (*pb.GetStoreResponse, error) {
	store, err := h.service.UpdateStore(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update store: %v", err)
	}
	return &pb.GetStoreResponse{Store: toProto(store)}, nil
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

func (h *StoreGRPCHandler) GetStoreByCode(ctx context.Context, req *pb.GetStoreByCodeRequest) (*pb.GetStoreResponse, error) {
	store, err := h.service.GetStoreByCode(ctx, req.StoreCode)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "store not found: %v", err)
	}
	return &pb.GetStoreResponse{Store: toProto(store)}, nil
}
