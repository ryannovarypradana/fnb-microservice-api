// internal/order/repository.go
package order

import (
	"fnb-system/pkg/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderInTx(order *model.Order, items *[]model.OrderItem) (*model.Order, error)
	FindOrdersByUserID(userID uint) (*[]model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) FindOrdersByUserID(userID uint) (*[]model.Order, error) {
	var orders []model.Order
	// Ambil semua order beserta item-item dan menu terkaitnya
	err := r.db.Preload("OrderItems.Menu").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

func (r *orderRepository) CreateOrderInTx(order *model.Order, items *[]model.OrderItem) (*model.Order, error) {
	// Mulai transaksi
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 1. Simpan data order utama
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback() // Batalkan transaksi jika gagal
		return nil, err
	}

	// 2. Kaitkan setiap item dengan OrderID yang baru dibuat, lalu simpan
	for i := range *items {
		(*items)[i].OrderID = order.ID
	}
	if err := tx.Create(items).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Jika semua berhasil, commit transaksi
	return order, tx.Commit().Error
}
