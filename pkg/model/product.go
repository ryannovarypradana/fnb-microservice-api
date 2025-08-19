package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	StoreID     uuid.UUID      `json:"store_id" gorm:"type:uuid"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}
