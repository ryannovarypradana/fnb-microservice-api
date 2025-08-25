// File: ryannovarypradana/fnb-microservice-api/fnb-microservice-api-dd6285232082f71efc6950ba298fd97bc68fbcc3/internal/user/repository.go

package user

import (
	"context"
	"errors"

	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *pb.GetAllUsersRequest) ([]*model.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// Fungsi baru untuk mengambil semua user dengan filter dan paginasi
func (r *userRepository) FindAll(ctx context.Context, req *pb.GetAllUsersRequest) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	if req.StoreId != nil {
		query = query.Where("store_id = ?", *req.StoreId)
	}
	if req.CompanyId != nil {
		query = query.Where("company_id = ?", *req.CompanyId)
	}
	if req.Search != nil {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+*req.Search+"%", "%"+*req.Search+"%")
	}

	// Hitung total data sebelum paginasi
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Terapkan paginasi
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(int(offset)).Limit(int(req.Limit)).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
