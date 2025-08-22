// pkg/model/jwt.go

package model

import "github.com/golang-jwt/jwt/v5"

// Claims adalah custom claims untuk JWT.
type Claims struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	CompanyID string `json:"companyID"`
	jwt.RegisteredClaims
}
