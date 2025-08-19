package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, store *model.Store) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Store, error)
	FindAll(ctx context.Context, search string) ([]*model.Store, error)
	FindByCode(ctx context.Context, code string) (*model.Store, error)
}

type storeRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(ctx context.Context, store *model.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *storeRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Store, error) {
	var store model.Store
	err := r.db.WithContext(ctx).Preload("Company").First(&store, "id = ?", id).Error
	return &store, err
}

func (r *storeRepository) FindAll(ctx context.Context, search string) ([]*model.Store, error) {
	var stores []*model.Store
	query := r.db.WithContext(ctx)
	if search != "" {
		query = query.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", search))
	}
	err := query.Find(&stores).Error
	return stores, err
}

func (r *storeRepository) FindByCode(ctx context.Context, code string) (*model.Store, error) {
	var store model.Store
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&store).Error
	return &store, err
}
