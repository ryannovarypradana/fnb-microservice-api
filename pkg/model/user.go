// pkg/model/user.go
package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Name     string    `json:"name" gorm:"type:varchar(255);not null"`
	Email    string    `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password string    `json:"-" gorm:"type:varchar(255);not null"`
	Role     string    `json:"role" gorm:"type:varchar(50);not null"`

	// Relasi untuk multi-tenancy
	CompanyID *uuid.UUID `json:"company_id,omitempty" gorm:"type:uuid"`
	StoreID   *uuid.UUID `json:"store_id,omitempty" gorm:"type:uuid"`

	// Associations (opsional, untuk kemudahan query)
	Company *Company `json:"-" gorm:"foreignkey:CompanyID"`
	Store   *Store   `json:"-" gorm:"foreignkey:StoreID"`

	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	// Default role jika tidak diset
	if u.Role == "" {
		u.Role = RoleCustomer
	}
	return
}
