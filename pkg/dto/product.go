// pkg/dto/product.go
package dto

// CategoryRequest adalah struct untuk membuat atau memperbarui kategori.
type CategoryRequest struct {
	Name string `json:"name"`
}

// MenuRequest adalah struct untuk membuat atau memperbarui menu.
type MenuRequest struct {
	Name       string  `json:"name"`
	CategoryID uint    `json:"category_id"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
}
