package auth

import (
	"errors"
	"fnb-system/pkg/dto"
	"fnb-system/pkg/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// --- Interface (Kontrak) ---
type AuthService interface {
	Register(req dto.AuthRegisterRequest) (*model.User, error)
	Login(email, password string) (string, error)
}

// Dependensi yang dibutuhkan oleh AuthService
type AuthRepository interface {
	CreateUser(user *model.User) (*model.User, error)
}
type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
}

// --- Implementation ---
type authService struct {
	authRepository AuthRepository
	userRepository UserRepository
}

// NewAuthService membuat instance baru dari authService.
func NewAuthService(authRepo AuthRepository, userRepo UserRepository) AuthService {
	return &authService{
		authRepository: authRepo,
		userRepository: userRepo,
	}
}

func (s *authService) Register(req dto.AuthRegisterRequest) (*model.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user", // Default role
	}

	createdUser, err := s.authRepository.CreateUser(&newUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Bandingkan password yang diinput dengan hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Jika berhasil, generate token
	return GenerateToken(user)
}

// --- JWT Helper Functions ---

// Claims mendefinisikan data yang akan disimpan di dalam token.
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken membuat token JWT baru untuk user.
func GenerateToken(user *model.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken memverifikasi token string dan mengembalikan claims jika valid.
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
