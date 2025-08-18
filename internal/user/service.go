// internal/user/service.go
package user

import (
	"errors"
	"fnb-system/pkg/dto"
	"fnb-system/pkg/model"
)

// --- Interface ---
// Ini adalah "kontrak" atau "perjanjian".
// Service hanya peduli pada fungsi-fungsi ini, bukan bagaimana cara kerjanya.
type UserRepository interface {
	FindAll(pagination *dto.Pagination) (*[]model.User, error)
	FindByID(id uint) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	Delete(id uint) error
}

// UserService mendefinisikan kontrak untuk logika bisnis user.
type UserService interface {
	GetAll(pagination *dto.Pagination) (*[]model.User, error)
	GetByID(id uint) (*model.User, error)
	Update(id uint, req dto.UserUpdateRequest) (*model.User, error)
	Delete(id uint) error
}

// --- Implementation ---
type userService struct {
	userRepository UserRepository
}

// NewUserService membuat instance baru dari userService.
// Perhatikan bahwa ia menerima INTERFACE, bukan struct konkret.
func NewUserService(userRepository UserRepository) UserService {
	return &userService{userRepository}
}

func (s *userService) GetAll(pagination *dto.Pagination) (*[]model.User, error) {
	return s.userRepository.FindAll(pagination)
}

func (s *userService) GetByID(id uint) (*model.User, error) {
	return s.userRepository.FindByID(id)
}

func (s *userService) Update(id uint, req dto.UserUpdateRequest) (*model.User, error) {
	// 1. Cari user yang ada
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Terapkan perubahan dari DTO ke model
	user.Name = req.Name
	user.Email = req.Email

	// 3. Simpan perubahan
	updatedUser, err := s.userRepository.Update(user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) Delete(id uint) error {
	// Cek dulu apakah user ada
	_, err := s.userRepository.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepository.Delete(id)
}
