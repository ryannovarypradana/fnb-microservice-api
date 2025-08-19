// /internal/auth/handler.go
package auth

import (
	"context"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type GRPCHandler struct {
	auth.UnimplementedAuthServiceServer
	service Service
}

func NewGRPCHandler(s Service) *GRPCHandler {
	return &GRPCHandler{service: s}
}

func (h *GRPCHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	token, err := h.service.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &auth.LoginResponse{Token: token}, nil
}

func (h *GRPCHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	user := &model.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	createdUser, err := h.service.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		User: &auth.User{
			Id:    createdUser.ID.String(),
			Name:  createdUser.Name,
			Email: createdUser.Email,
			Role:  createdUser.Role,
		},
	}, nil
}
