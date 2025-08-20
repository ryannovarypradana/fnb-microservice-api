// /internal/auth/service.go
package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, user *model.User) (*model.User, error)
}

type authService struct {
	repo Repository
	db   *gorm.DB
}

func NewService(repo Repository, db *gorm.DB) Service {
	return &authService{repo: repo, db: db}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
	// Sisipkan ID perusahaan & toko jika ada
	if user.CompanyID != nil {
		claims["company_id"] = user.CompanyID.String()
	}
	if user.StoreID != nil {
		claims["store_id"] = user.StoreID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) Register(ctx context.Context, user *model.User) (*model.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	// Set default role jika kosong
	if user.Role == "" {
		user.Role = model.RoleCustomer
	}

	// Create user
	if err := s.repo.CreateUser(ctx, nil, user); err != nil {
		return nil, err
	}

	// Menghapus password dari response
	user.Password = ""
	return user, nil
}
