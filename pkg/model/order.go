package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;"`
	UserID     uuid.UUID      `json:"user_id" gorm:"type:uuid"`
	Status     string         `json:"status"`
	TotalPrice float64        `json:"total_price"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// FIX: Add the relationship to OrderItem
	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	OrderID uuid.UUID `json:"order_id" gorm:"type:uuid"`
	// FIX: Field should be ProductID to match the service logic
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid"`
	Quantity  int       `json:"quantity"`
	// FIX: Field should be Subtotal
	Subtotal  float64        `json:"subtotal"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
