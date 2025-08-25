// File: ryannovarypradana/fnb-microservice-api/fnb-microservice-api-dd6285232082f71efc6950ba298fd97bc68fbcc3/pkg/model/company.go

package model

import (
	"time" // <-- TAMBAHKAN IMPORT

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Company struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name    string    `json:"name" gorm:"type:varchar(255);not null;unique"`
	Address string    `json:"address"`
	Code    string    `json:"code" gorm:"type:varchar(10);not null;unique"`

	// --- TAMBAHKAN FIELD INI ---
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Company) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return
}
