package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OpeningHours defines the structure for operational hours.
type OpeningHours struct {
	Monday    string `json:"monday"`
	Tuesday   string `json:"tuesday"`
	Wednesday string `json:"wednesday"`
	Thursday  string `json:"thursday"`
	Friday    string `json:"friday"`
	Saturday  string `json:"saturday"`
	Sunday    string `json:"sunday"`
}

// GORM support for OpeningHours struct
func (o *OpeningHours) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &o)
}

func (o OpeningHours) Value() (driver.Value, error) {
	return json.Marshal(o)
}

// Store defines the GORM model for a store.
type Store struct {
	ID               uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code             string        `gorm:"type:varchar(10);unique;not null" json:"code"`
	Name             string        `gorm:"type:varchar(255);not null" json:"name"`
	Location         string        `json:"location"` // Renamed from Address
	TaxPercentage    *float64      `gorm:"not null;default:11.0" json:"taxPercentage"`
	CompanyID        uuid.UUID     `gorm:"type:uuid;not null" json:"companyId"`
	OperationalHours *OpeningHours `gorm:"type:jsonb" json:"operationalHours"`
	Latitude         *float64      `json:"latitude"`
	Longitude        *float64      `json:"longitude"`
	BannerImageURL   string        `json:"bannerImageUrl"`

	Company Company `gorm:"foreignKey:CompanyID" json:"-"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Store) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	// Add logic to generate a unique 'Code' if it's empty
	return
}
