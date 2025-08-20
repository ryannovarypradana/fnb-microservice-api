package model

import (
	"gorm.io/gorm"
)

// Category merepresentasikan kategori dari sebuah menu di dalam toko.
type Category struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	StoreID uint   `gorm:"not null"`
	Store   Store
	Menus   []Menu // Sebuah kategori bisa memiliki banyak menu
}

// Menu merepresentasikan item yang bisa dijual (sebelumnya Product).
type Menu struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null"`
	Description string
	Price       float64 `gorm:"not null;default:0"`
	ImageURL    string
	StoreID     uint `gorm:"not null"`
	Store       Store
	CategoryID  uint `gorm:"not null"`
	Category    Category
}
