package model

import (
	"time"

	"gorm.io/gorm"
)

// Bundle represents a collection of products sold together at a special price.
type Bundle struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string  `gorm:"not null;size:255"`
	Price float64 `gorm:"not null"`

	// This field will hold the JSON string representation of ProductIDs in the database.
	Products string `gorm:"type:text"`

	// This field is for application logic and will not be stored directly in the database.
	// We use `gorm:"-"` to ignore this field during database operations.
	ProductIDs []uint `gorm:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
