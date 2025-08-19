package user

import (
	"context"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
)

type GRPCHandler struct {
	user.UnimplementedUserServiceServer // Penting untuk forward compatibility
	service                             Service
}

func NewGRPCHandler(s Service) *GRPCHandler {
	return &GRPCHandler{service: s}
}

// GetUser adalah implementasi dari RPC GetUser yang ada di user.proto
func (h *GRPCHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	foundUser, err := h.service.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	// Terjemahkan model internal ke response message protobuf
	res := &user.GetUserResponse{
		User: &user.User{
			Id:    foundUser.ID.String(),
			Name:  foundUser.Name,
			Email: foundUser.Email,
			Role:  foundUser.Role,
		},
	}

	if foundUser.CompanyID != nil {
		res.User.CompanyId = foundUser.CompanyID.String()
	}
	if foundUser.StoreID != nil {
		res.User.StoreId = foundUser.StoreID.String()
	}

	return res, nil
}

// Implementasikan RPC lain dari user.proto di sini...
