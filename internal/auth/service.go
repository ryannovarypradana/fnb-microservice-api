package auth

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"

	// FIXED: Typo 'github.comcom' corrected to 'github.com'
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/utils"
)

// AuthService is the interface for authentication business logic.
type AuthService interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*model.User, error)
	RegisterStaff(ctx context.Context, req *pb.RegisterStaffRequest) (*model.User, error)
	Login(ctx context.Context, req *pb.LoginRequest) (string, error)
}

type authService struct {
	repo       AuthRepository
	jwtService utils.JwtService
}

// NewAuthService is the constructor for authService.
func NewAuthService(repo AuthRepository, jwtService utils.JwtService) AuthService {
	return &authService{repo, jwtService}
}

// Register is for customer registration.
func (s *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*model.User, error) {
	staffReq := &pb.RegisterStaffRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     string(model.RoleCustomer),
	}
	return s.RegisterStaff(ctx, staffReq)
}

// RegisterStaff is for internal staff registration.
func (s *authService) RegisterStaff(ctx context.Context, req *pb.RegisterStaffRequest) (*model.User, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR: failed to hash password: %v", err)
		return nil, errors.New("failed to process request")
	}

	newUser := &model.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     model.Role(req.Role),
	}

	if req.CompanyId != "" {
		companyID, err := uuid.Parse(req.CompanyId)
		if err == nil {
			newUser.CompanyID = &companyID
		}
	}

	if req.StoreId != "" {
		storeID, err := uuid.Parse(req.StoreId)
		if err == nil {
			newUser.StoreID = &storeID
		}
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		log.Printf("ERROR: failed to register user: %v", err)
		return nil, errors.New("failed to register user")
	}

	return newUser, nil
}

func (s *authService) Login(ctx context.Context, req *pb.LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
