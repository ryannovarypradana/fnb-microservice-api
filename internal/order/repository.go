package order

import (
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

// IRepository adalah interface untuk order repository
type IRepository interface {
	CreateOrderWithItems(order *model.Order, items []*model.OrderItem) error
	FindOrderByID(orderID uint) (*model.Order, error)
	FindOrdersByUserID(userID uint) ([]*model.Order, error)
}

type Repository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IRepository {
	return &Repository{db: db}
}

func (r *Repository) CreateOrderWithItems(order *model.Order, items []*model.OrderItem) error {
	// Memulai transaksi
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan order utama untuk mendapatkan ID
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Set OrderID untuk setiap item dan simpan
		for _, item := range items {
			item.OrderID = order.ID
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repository) FindOrderByID(orderID uint) (*model.Order, error) {
	var order model.Order
	// Preload OrderItems untuk mendapatkan detail item
	if err := r.db.Preload("OrderItems").First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) FindOrdersByUserID(userID uint) ([]*model.Order, error) {
	var orders []*model.Order
	if err := r.db.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
