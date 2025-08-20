package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;"`
	CompanyID *uuid.UUID `json:"company_id,omitempty" gorm:"type:uuid"`
	StoreID   *uuid.UUID `json:"store_id,omitempty" gorm:"type:uuid"`
	Name      string     `json:"name" gorm:"type:varchar(255);not null"`
	Email     string     `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password  string     `json:"-" gorm:"not null"`
	// CORRECTED: The type is now model.Role, not string.
	Role Role `json:"role" gorm:"type:varchar(50);not null"`

	// Relasi
	Company Company `json:"-" gorm:"foreignKey:CompanyID"`
	Store   Store   `json:"-" gorm:"foreignKey:StoreID"`

	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	// Set default role if not provided
	if u.Role == "" {
		// This assignment is now valid because both sides are of type model.Role
		u.Role = RoleCustomer
	}
	return
}
