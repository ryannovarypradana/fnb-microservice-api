package dto

// CreateStoreRequest adalah DTO untuk membuat toko baru.
type CreateStoreRequest struct {
	Name      string `json:"name" binding:"required"`
	Address   string `json:"address" binding:"required"`
	CompanyID uint   `json:"company_id" binding:"required"`
}

// UpdateStoreRequest adalah DTO untuk memperbarui data toko.
type UpdateStoreRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}
