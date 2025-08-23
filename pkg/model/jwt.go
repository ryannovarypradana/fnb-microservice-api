// pkg/model/jwt.go

package model

import "github.com/golang-jwt/jwt/v5"

// Claims adalah custom claims untuk JWT.
type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CompanyID string `json:"company_id"`
	StoreID   string `json:"store_id"`
	jwt.RegisteredClaims
}
