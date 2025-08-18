// pkg/dto/auth.go
package dto

// AuthRegisterRequest adalah struct untuk request body saat registrasi.
type AuthRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthLoginRequest adalah struct untuk request body saat login.
type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
