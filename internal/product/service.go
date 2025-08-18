// internal/product/service.go
package product

import (
	"fnb-system/pkg/dto"
	"fnb-system/pkg/model"
)

// Tambahkan GetMenuByID ke interface
type ProductService interface {
	CreateCategory(req dto.CategoryRequest) (*model.Category, error)
	GetAllCategories() (*[]model.Category, error)
	CreateMenu(req dto.MenuRequest) (*model.Menu, error)
	GetAllMenus() (*[]model.Menu, error)
	GetMenuByID(id uint) (*model.Menu, error) // <-- Tambahkan ini
}

type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo}
}

// Tambahkan implementasi fungsi GetMenuByID
func (s *productService) GetMenuByID(id uint) (*model.Menu, error) {
	return s.repo.FindMenuByID(id)
}

// --- Fungsi yang sudah ada sebelumnya ---

func (s *productService) CreateCategory(req dto.CategoryRequest) (*model.Category, error) {
	category := model.Category{
		Name: req.Name,
	}
	return s.repo.CreateCategory(&category)
}

func (s *productService) GetAllCategories() (*[]model.Category, error) {
	return s.repo.FindAllCategories()
}

func (s *productService) CreateMenu(req dto.MenuRequest) (*model.Menu, error) {
	menu := model.Menu{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Price:      req.Price,
		ImageURL:   req.ImageURL,
	}
	return s.repo.CreateMenu(&menu)
}

func (s *productService) GetAllMenus() (*[]model.Menu, error) {
	return s.repo.FindAllMenus() // <-- INI YANG BENAR
}
