package store

import (
	"context"
	"errors"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type StoreRepository interface {
	Create(ctx context.Context, store *model.Store) error
	FindByID(ctx context.Context, id string) (*model.Store, error)
	FindAll(ctx context.Context, search string) ([]*model.Store, error)
	Update(ctx context.Context, store *model.Store) error
	Delete(ctx context.Context, id string) error
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(ctx context.Context, store *model.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *storeRepository) FindByID(ctx context.Context, id string) (*model.Store, error) {
	var store model.Store
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&store).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("store not found")
		}
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) FindAll(ctx context.Context, search string) ([]*model.Store, error) {
	var stores []*model.Store
	query := r.db.WithContext(ctx)
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := query.Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *storeRepository) Update(ctx context.Context, store *model.Store) error {
	return r.db.WithContext(ctx).Save(store).Error
}

func (r *storeRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.Store{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("store not found")
	}
	return nil
}
