package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

// JwtService adalah interface untuk operasi JWT.
type JwtService interface {
	GenerateToken(userID uuid.UUID, role string) (string, error)
	VerifyToken(tokenString string) (*model.Claims, error)
}

// jwtService adalah implementasi dari JwtService.
type jwtService struct {
	secret string
}

// NewJwtService adalah konstruktor untuk jwtService.
func NewJwtService(cfg *config.Config) JwtService {
	return &jwtService{
		secret: cfg.App.JWTSecret,
	}
}

// GenerateToken sekarang adalah sebuah method.
func (s *jwtService) GenerateToken(userID uuid.UUID, role string) (string, error) {
	claims := model.Claims{
		UserID: userID.String(), // Menggunakan String() dari UUID
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// VerifyToken sekarang adalah sebuah method.
func (s *jwtService) VerifyToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
