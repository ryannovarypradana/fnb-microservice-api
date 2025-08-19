package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

// FIX: Tambahkan metode `Create` ke dalam interface.
type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Create(ctx context.Context, user *model.User) error // <-- TAMBAHKAN BARIS INI
}

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}

// FIX: Tambahkan implementasi dari metode `Create`.
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	// Fungsi ini akan membuat record user baru di dalam database.
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
