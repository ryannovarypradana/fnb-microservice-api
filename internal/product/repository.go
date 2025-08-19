package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, product *model.Product) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	FindAll(ctx context.Context, search, storeID string) ([]*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(ctx context.Context, search, storeID string) ([]*model.Product, error) {
	var products []*model.Product
	query := r.db.WithContext(ctx)
	if search != "" {
		query = query.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", search))
	}
	if storeID != "" {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, "id = ?", id).Error
}
