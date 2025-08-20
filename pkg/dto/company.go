// pkg/dto/company.go

package dto

// CreateCompanyRequest adalah DTO untuk membuat perusahaan baru.
type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}
