package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name    string    `json:"name" gorm:"type:varchar(255);not null;unique"`
	Address string    `json:"address"`
	// PASTIKAN FIELD INI ADA
	Code string `json:"code" gorm:"type:varchar(10);not null;unique"`

	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Company) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
