// pkg/model/product.go
package model

import "time"

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Menus     []Menu    `json:"menus"` // Relasi one-to-many
}

type Menu struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	CategoryID uint      `json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Price      float64   `gorm:"not null" json:"price"`
	ImageURL   string    `gorm:"size:255" json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
