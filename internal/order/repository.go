// /internal/order/repository.go
package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, order *model.Order, items []model.OrderItem) error // FIX: Changed type here
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order, items []model.OrderItem) error { // FIX: Changed type here
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = order.ID
		}
		if err := tx.Create(&items).Error; err != nil { // Pass address of slice
			return err
		}
		return nil
	})
}

// ... rest of the file is the same
func (r *orderRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	var order model.Order
	if err := r.db.WithContext(ctx).Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	var orders []*model.Order
	if err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
