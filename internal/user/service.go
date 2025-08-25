// File: ryannovarypradana/fnb-microservice-api/fnb-microservice-api-dd6285232082f71efc6950ba298fd97bc68fbcc3/internal/user/service.go

package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	companyPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
	userPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type UserService interface {
	GetUser(ctx context.Context, req *userPB.GetUserRequest) (*model.User, error)
	RegisterStaff(ctx context.Context, req *userPB.RegisterStaffRequest) (*model.User, error)
	UpdateUser(ctx context.Context, req *userPB.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
	CreateCompanyWithRep(ctx context.Context, req *userPB.CreateCompanyWithRepRequest) (*companyPB.Company, *model.User, error)
	GetAllUsers(ctx context.Context, req *userPB.GetAllUsersRequest) ([]*model.User, int64, error)
}

type userService struct {
	repo          UserRepository
	companyClient companyPB.CompanyServiceClient
}

func NewUserService(repo UserRepository, companyClient companyPB.CompanyServiceClient) UserService {
	return &userService{repo: repo, companyClient: companyClient}
}

func (s *userService) GetUser(ctx context.Context, req *userPB.GetUserRequest) (*model.User, error) {
	return s.repo.FindByID(ctx, req.Id)
}

// Fungsi baru untuk memanggil repository
func (s *userService) GetAllUsers(ctx context.Context, req *userPB.GetAllUsersRequest) ([]*model.User, int64, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *userService) RegisterStaff(ctx context.Context, req *userPB.RegisterStaffRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	storeID, _ := uuid.Parse(req.StoreId)

	newUser := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     model.Role(req.Role),
		StoreID:  &storeID,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *userPB.UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) CreateCompanyWithRep(ctx context.Context, req *userPB.CreateCompanyWithRepRequest) (*companyPB.Company, *model.User, error) {
	createCompanyReq := &companyPB.CreateCompanyRequest{
		Name:    req.CompanyName,
		Address: req.CompanyAddress,
	}

	companyRes, err := s.companyClient.CreateCompany(ctx, createCompanyReq)
	if err != nil {
		return nil, nil, errors.New("failed to create company via gRPC: " + err.Error())
	}
	newCompany := companyRes.Company

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, errors.New("failed to hash password")
	}

	companyID, _ := uuid.Parse(newCompany.Id)
	newUser := &model.User{
		Name:      req.UserName,
		Email:     req.UserEmail,
		Password:  string(hashedPassword),
		Role:      model.RoleCompanyRep,
		CompanyID: &companyID,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, nil, errors.New("failed to create company representative user: " + err.Error())
	}

	return newCompany, newUser, nil
}
