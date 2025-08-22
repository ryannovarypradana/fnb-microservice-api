package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category merepresentasikan kategori dari sebuah menu di dalam toko.
type Category struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name    string    `gorm:"type:varchar(100);not null"`
	StoreID uuid.UUID `gorm:"type:uuid;not null"`
	Store   Store
	Menus   []Menu // Sebuah kategori bisa memiliki banyak menu

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Menu merepresentasikan item yang bisa dijual (sebelumnya Product).
type Menu struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string
	Price       float64 `gorm:"not null;default:0"`
	ImageURL    string
	StoreID     uuid.UUID `gorm:"type:uuid;not null"`
	Store       Store
	CategoryID  uuid.UUID `gorm:"type:uuid;not null"`
	Category    Category

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Hooks GORM untuk model Category
func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}

// Hooks GORM untuk model Menu
func (m *Menu) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
