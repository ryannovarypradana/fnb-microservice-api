package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"

	// CORRECTED: Added the missing import for the utils package
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/utils"
)

// AuthService is the interface for authentication business logic.
type AuthService interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req *pb.LoginRequest) (string, error)
}

type authService struct {
	repo       AuthRepository
	jwtService utils.JwtService // This line now works because of the import
}

// NewAuthService is the constructor for authService.
func NewAuthService(repo AuthRepository, jwtService utils.JwtService) AuthService {
	return &authService{repo, jwtService}
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

func (s *authService) Register(ctx context.Context, req *pb.RegisterRequest) (*model.User, error) {
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	newUser := &model.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     model.RoleCustomer,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, errors.New("failed to register user")
	}

	return newUser, nil
}
