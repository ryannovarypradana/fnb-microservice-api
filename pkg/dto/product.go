package dto

// CreateMenuRequest defines the structure for creating a new menu.
type CreateMenuRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
	CategoryID  string  `json:"category_id" binding:"required"`
	StoreID     string  `json:"store_id" binding:"required"`
}

// UpdateMenuRequest defines the structure for updating an existing menu.
type UpdateMenuRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
	CategoryID  string  `json:"category_id" binding:"required"`
}

// CreateCategoryRequest defines the structure for creating a new category.
type CreateCategoryRequest struct {
	Name    string `json:"name" binding:"required"`
	StoreID string `json:"store_id" binding:"required"`
}

// UpdateCategoryRequest defines the structure for updating an existing category.
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// MenuResponse defines the standard structure for a menu response.
type MenuResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	ImageURL     string  `json:"image_url"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	StoreID      string  `json:"store_id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// CategoryResponse defines the standard structure for a category response.
type CategoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StoreID   string `json:"store_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
