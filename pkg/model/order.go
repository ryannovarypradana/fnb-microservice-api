// pkg/model/order.go
package model

import "time"

type Order struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	UserID     uint        `json:"user_id"`
	User       User        `gorm:"foreignKey:UserID" json:"user"`
	TotalPrice float64     `gorm:"not null" json:"total_price"`
	Status     string      `gorm:"size:50;not null;default:'pending'" json:"status"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	OrderID  uint    `json:"order_id"`
	MenuID   uint    `json:"menu_id"`
	Menu     Menu    `gorm:"foreignKey:MenuID" json:"menu"`
	Quantity int     `gorm:"not null" json:"quantity"`
	Price    float64 `gorm:"not null" json:"price"` // Harga saat item dipesan
}
