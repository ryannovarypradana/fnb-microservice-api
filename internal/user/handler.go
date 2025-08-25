// File: ryannovarypradana/fnb-microservice-api/fnb-microservice-api-dd6285232082f71efc6950ba298fd97bc68fbcc3/internal/user/handler.go
package user

import (
	"context"

	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserGRPCHandler struct {
	pb.UnimplementedUserServiceServer
	service UserService
}

func NewUserGRPCHandler(grpcServer *grpc.Server, service UserService) {
	handler := &UserGRPCHandler{service: service}
	pb.RegisterUserServiceServer(grpcServer, handler)
}

// Fungsi baru untuk menangani RPC GetAllUsers
func (h *UserGRPCHandler) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, total, err := h.service.GetAllUsers(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUser := &pb.User{
			Id:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  string(user.Role),
		}
		if user.CompanyID != nil {
			companyIDStr := user.CompanyID.String()
			pbUser.CompanyId = &companyIDStr
		}
		if user.StoreID != nil {
			storeIDStr := user.StoreID.String()
			pbUser.StoreId = &storeIDStr
		}
		pbUsers = append(pbUsers, pbUser)
	}

	return &pb.GetAllUsersResponse{
		Users: pbUsers,
		Total: int32(total),
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (h *UserGRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	foundUser, err := h.service.GetUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	res := &pb.GetUserResponse{
		User: &pb.User{
			Id:    foundUser.ID.String(),
			Name:  foundUser.Name,
			Email: foundUser.Email,
			Role:  string(foundUser.Role),
		},
	}

	if foundUser.CompanyID != nil {
		companyIDStr := foundUser.CompanyID.String()
		res.User.CompanyId = &companyIDStr
	}
	if foundUser.StoreID != nil {
		storeIDStr := foundUser.StoreID.String()
		res.User.StoreId = &storeIDStr
	}

	return res, nil
}

func (h *UserGRPCHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	updatedUser, err := h.service.UpdateUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}
	return h.GetUser(ctx, &pb.GetUserRequest{Id: updatedUser.ID.String()})
}

func (h *UserGRPCHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if err := h.service.DeleteUser(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return &pb.DeleteUserResponse{Message: "User deleted successfully"}, nil
}

func (h *UserGRPCHandler) CreateCompanyWithRep(ctx context.Context, req *pb.CreateCompanyWithRepRequest) (*pb.CreateCompanyWithRepResponse, error) {
	company, user, err := h.service.CreateCompanyWithRep(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create company with rep: %v", err)
	}

	resUser := &pb.User{
		Id:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}

	if user.CompanyID != nil {
		companyIDStr := user.CompanyID.String()
		resUser.CompanyId = &companyIDStr
	}

	return &pb.CreateCompanyWithRepResponse{
		CompanyId:   company.Id,
		CompanyName: company.Name,
		User:        resUser,
	}, nil
}
