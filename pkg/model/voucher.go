package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Voucher represents a code that can be redeemed for a discount or benefit.
type Voucher struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code      string    `gorm:"unique;not null;size:100;index"`
	Type      string    `gorm:"not null;size:50"` // e.g., "PERCENTAGE", "FIXED_AMOUNT"
	Value     float64   `gorm:"not null"`
	Quota     int       `gorm:"not null"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Hooks GORM untuk model Voucher
func (v *Voucher) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return
}
