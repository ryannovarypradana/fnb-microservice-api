package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	CompanyID uuid.UUID `json:"company_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Address   string    `json:"address"`
	// PASTIKAN FIELD INI ADA
	Code string `json:"code" gorm:"type:varchar(10);not null;unique"`

	// Relasi
	Company Company `json:"-" gorm:"foreignKey:CompanyID"`

	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (s *Store) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
