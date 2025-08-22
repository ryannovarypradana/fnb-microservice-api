package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Bundle represents a collection of products sold together at a special price.
type Bundle struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name  string    `gorm:"not null;size:255"`
	Price float64   `gorm:"not null"`

	// This field will hold the JSON string representation of ProductIDs in the database.
	Products string `gorm:"type:text"`

	// This field is for application logic and will not be stored directly in the database.
	// We use `gorm:"-"` to ignore this field during database operations.
	ProductIDs []uuid.UUID `gorm:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Hooks GORM untuk model Bundle
func (b *Bundle) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
