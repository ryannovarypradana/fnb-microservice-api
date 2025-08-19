// /internal/auth/repository.go
package auth

import (
	"context"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, tx *gorm.DB, user *model.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUser(ctx context.Context, tx *gorm.DB, user *model.User) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(user).Error
}
