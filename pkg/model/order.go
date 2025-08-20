package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderStatus mendefinisikan tipe untuk status pesanan.
type OrderStatus string

// Konstanta untuk berbagai status pesanan.
const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusCompleted OrderStatus = "COMPLETED"
	StatusCanceled  OrderStatus = "CANCELED"
)

// Order adalah model GORM untuk tabel 'orders'.
type Order struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;"`
	UserID      *uuid.UUID  `json:"user_id,omitempty" gorm:"type:uuid"` // Pointer karena bisa jadi pesanan publik
	StoreID     uuid.UUID   `json:"store_id" gorm:"type:uuid;not null"`
	OrderCode   string      `json:"order_code" gorm:"type:varchar(20);unique;not null"`
	TotalAmount float64     `json:"total_amount" gorm:"not null"`
	Status      OrderStatus `json:"status" gorm:"type:varchar(50);not null"`

	// Informasi tambahan untuk pesanan publik (dine-in)
	CustomerName string `json:"customer_name,omitempty"`
	TableNumber  string `json:"table_number,omitempty"`

	// Relasi
	User  User         `json:"-" gorm:"foreignKey:UserID"`
	Store Store        `json:"-" gorm:"foreignKey:StoreID"`
	Items []*OrderItem `json:"items" gorm:"foreignKey:OrderID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// OrderItem adalah model GORM untuk tabel 'order_items'.
type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"` // Harga satuan saat pesanan dibuat

	// Relasi
	// CORRECTED: Mengubah 'Product' menjadi 'Menu' agar sesuai dengan model Anda.
	Menu Menu `json:"product" gorm:"foreignKey:ProductID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate adalah hook GORM untuk menghasilkan ID baru sebelum membuat entitas.
func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}

// BeforeCreate adalah hook GORM untuk menghasilkan ID baru sebelum membuat entitas.
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return
}
