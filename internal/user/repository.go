package user

import (
	"fnb-system/pkg/dto"
	"fnb-system/pkg/model"

	"gorm.io/gorm"
)

// UserRepository mendefinisikan kontrak untuk operasi database user.
// Ini adalah implementasi konkretnya.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru dari userRepository.
func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

// FindAll mengambil semua user dengan paginasi.
func (r *userRepository) FindAll(pagination *dto.Pagination) (*[]model.User, error) {
	var users []model.User
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := r.db.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Model(&model.User{}).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

// FindByID mencari user berdasarkan ID.
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail mencari user berdasarkan email.
// Fungsi ini juga penting untuk Auth Service.
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update memperbarui data user di database.
func (r *userRepository) Update(user *model.User) (*model.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Delete menghapus user dari database.
func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
