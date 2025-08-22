package dto

import "time"

// CreateDiscountRequest is a DTO for creating a new discount.
type CreateDiscountRequest struct {
	Name      string    `json:"name" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	Value     float64   `json:"value" binding:"required,gt=0"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
}

// UpdateDiscountRequest is a DTO for updating an existing discount.
type UpdateDiscountRequest struct {
	Name      string    `json:"name,omitempty"`
	Type      string    `json:"type,omitempty"`
	Value     float64   `json:"value,omitempty" binding:"omitempty,gt=0"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

// CreateVoucherRequest is a DTO for creating a new voucher.
type CreateVoucherRequest struct {
	Code      string    `json:"code" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	Value     float64   `json:"value" binding:"required,gt=0"`
	Quota     int       `json:"quota" binding:"required,gt=0"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
}

// UpdateVoucherRequest is a DTO for updating an existing voucher.
type UpdateVoucherRequest struct {
	Code      string    `json:"code,omitempty"`
	Type      string    `json:"type,omitempty"`
	Value     float64   `json:"value,omitempty" binding:"omitempty,gt=0"`
	Quota     int       `json:"quota,omitempty" binding:"omitempty,gt=0"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

// CreateBundleRequest is a DTO for creating a new product bundle.
type CreateBundleRequest struct {
	Name       string   `json:"name" binding:"required"`
	Price      float64  `json:"price" binding:"required,gt=0"`
	ProductIDs []string `json:"product_ids" binding:"required,min=1"`
}

// UpdateBundleRequest is a DTO for updating an existing product bundle.
type UpdateBundleRequest struct {
	Name       string   `json:"name,omitempty"`
	Price      float64  `json:"price,omitempty" binding:"omitempty,gt=0"`
	ProductIDs []string `json:"product_ids,omitempty" binding:"omitempty,min=1"`
}

// DiscountResponse represents a single discount in the response.
type DiscountResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// VoucherResponse represents a single voucher in the response.
type VoucherResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Quota     int       `json:"quota"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// BundleResponse represents a single product bundle in the response.
type BundleResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	ProductIDs []string `json:"product_ids"`
}
