package auth

import (
	"context"

	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	service AuthService
}

// NewAuthGRPCHandler is the constructor for the gRPC handler.
func NewAuthGRPCHandler(grpcServer *grpc.Server, service AuthService) {
	handler := &AuthGRPCHandler{service: service}
	pb.RegisterAuthServiceServer(grpcServer, handler)
}

// toProtoUser converts a user model to a protobuf user message.
func toProtoUser(user *model.User) *pb.User {
	return &pb.User{
		Id:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
}

func (h *AuthGRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	createdUser, err := h.service.Register(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "registration failed: %v", err)
	}
	return &pb.RegisterResponse{User: toProtoUser(createdUser)}, nil
}

func (h *AuthGRPCHandler) RegisterStaff(ctx context.Context, req *pb.RegisterStaffRequest) (*pb.RegisterResponse, error) {
	createdUser, err := h.service.RegisterStaff(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "staff registration failed: %v", err)
	}
	return &pb.RegisterResponse{User: toProtoUser(createdUser)}, nil
}

func (h *AuthGRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := h.service.Login(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}
	return &pb.LoginResponse{Token: token}, nil
}
