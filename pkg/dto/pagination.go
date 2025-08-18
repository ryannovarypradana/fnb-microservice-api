// pkg/dto/pagination.go
package dto

// Pagination adalah struct untuk mengatur data paginasi dari query parameter.
type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}
