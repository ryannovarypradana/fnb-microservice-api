package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/middleware"
)

type JWTUtil struct {
	secretKey string
}

func NewJWTUtil(secretKey string) *JWTUtil {
	return &JWTUtil{secretKey: secretKey}
}

func (j *JWTUtil) GenerateToken(userID, role, companyID string) (string, error) {
	claims := &middleware.JWTClaims{
		UserID:    userID,
		Role:      role,
		CompanyID: companyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTUtil) ValidateToken(tokenString string) (*middleware.JWTClaims, error) {
	claims := &middleware.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
