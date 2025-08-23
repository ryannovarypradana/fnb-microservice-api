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
	// Tambahkan parameter companyID
	GenerateToken(userID uuid.UUID, email string, companyID string, storeID string, role string) (string, error)
	VerifyToken(tokenString string) (*model.Claims, error)
}

// jwtService adalah implementasi dari JwtService.
type jwtService struct {
	secret string
	// Tambahkan TTL untuk durasi token
	ttl time.Duration
}

// NewJwtService adalah konstruktor untuk jwtService.
func NewJwtService(cfg *config.Config) JwtService {
	return &jwtService{
		secret: cfg.App.JWTSecret,
		// Atur TTL dari config, default 24 jam jika tidak ada
		ttl: 24 * time.Hour,
	}
}

// GenerateToken sekarang menerima companyID.
func (s *jwtService) GenerateToken(userID uuid.UUID, email string, companyID string, storeID string, role string) (string, error) {
	claims := model.Claims{
		UserID:    userID.String(),
		Email:     email,
		CompanyID: companyID, // Tambahkan companyID ke dalam claims
		StoreID:   storeID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
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
