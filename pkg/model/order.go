package model

import (
	"gorm.io/gorm"
)

// Order merepresentasikan pesanan utama.
type Order struct {
	gorm.Model
	UserID      uint `gorm:"not null"`
	User        User
	StoreID     uint `gorm:"not null"`
	Store       Store
	TotalAmount float64 `gorm:"not null"`
	Status      string  `gorm:"type:varchar(50);not null;default:'PENDING'"`
	OrderItems  []OrderItem
}

// OrderItem merepresentasikan setiap item di dalam sebuah pesanan.
type OrderItem struct {
	gorm.Model
	OrderID uint `gorm:"not null"`
	// Perbaikan di sini: Menggunakan ProductID agar konsisten dengan relasi GORM
	ProductID uint    `gorm:"not null"`
	Menu      Menu    `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"` // Harga per item saat transaksi
}
