package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Discount represents a promotion that reduces the price of an item or total order.
type Discount struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `gorm:"not null;size:255"`
	Type      string    `gorm:"not null;size:50"` // e.g., "PERCENTAGE", "FIXED_AMOUNT"
	Value     float64   `gorm:"not null"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Hooks GORM untuk model Discount
func (d *Discount) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return
}
