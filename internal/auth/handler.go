package auth

import (
	"context"

	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	service AuthService // This type is now in the same package
}

// NewAuthGRPCHandler is the constructor for the gRPC handler.
func NewAuthGRPCHandler(grpcServer *grpc.Server, service AuthService) {
	handler := &AuthGRPCHandler{service: service}
	pb.RegisterAuthServiceServer(grpcServer, handler)
}

func (h *AuthGRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := h.service.Login(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}
	return &pb.LoginResponse{Token: token}, nil
}

func (h *AuthGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	createdUser, err := h.service.Register(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "registration failed: %v", err)
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Id:    createdUser.ID.String(),
			Name:  createdUser.Name,
			Email: createdUser.Email,
			Role:  string(createdUser.Role),
		},
	}, nil
}
