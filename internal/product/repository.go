package product

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

// IRepository mendefinisikan antarmuka untuk operasi data produk.
type IRepository interface {
	// Operasi Menu
	CreateMenu(menu *model.Menu) error
	GetMenuByID(menuID uuid.UUID) (*model.Menu, error)
	UpdateMenu(menu *model.Menu) error
	DeleteMenu(menuID uuid.UUID) error
	FindMenusByStoreID(storeID uuid.UUID) ([]*model.Menu, error)

	// Operasi Kategori
	CreateCategory(category *model.Category) error
	GetCategoryByID(categoryID uuid.UUID) (*model.Category, error)
	UpdateCategory(category *model.Category) error
	DeleteCategory(categoryID uuid.UUID) error
	FindCategoriesByStoreID(storeID uuid.UUID) ([]*model.Category, error)
}

// Repository adalah implementasi GORM dari IRepository.
type Repository struct {
	db *gorm.DB
}

// NewProductRepository adalah konstruktor untuk Repository.
func NewProductRepository(db *gorm.DB) IRepository {
	return &Repository{db: db}
}

// --- Implementasi Menu ---

func (r *Repository) CreateMenu(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *Repository) GetMenuByID(menuID uuid.UUID) (*model.Menu, error) {
	var menu model.Menu
	if err := r.db.Preload("Category").First(&menu, "id = ?", menuID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("menu not found")
		}
		return nil, err
	}
	return &menu, nil
}

func (r *Repository) UpdateMenu(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

func (r *Repository) DeleteMenu(menuID uuid.UUID) error {
	result := r.db.Delete(&model.Menu{}, "id = ?", menuID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("menu not found")
	}
	return nil
}

func (r *Repository) FindMenusByStoreID(storeID uuid.UUID) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.db.Preload("Category").Where("store_id = ?", storeID).Find(&menus).Error
	return menus, err
}

// --- Implementasi Kategori ---

func (r *Repository) CreateCategory(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *Repository) GetCategoryByID(categoryID uuid.UUID) (*model.Category, error) {
	var category model.Category
	if err := r.db.First(&category, "id = ?", categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (r *Repository) UpdateCategory(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *Repository) DeleteCategory(categoryID uuid.UUID) error {
	// Pastikan tidak ada menu yang terhubung sebelum menghapus
	var count int64
	r.db.Model(&model.Menu{}).Where("category_id = ?", categoryID).Count(&count)
	if count > 0 {
		return errors.New("category has associated menus and cannot be deleted")
	}
	result := r.db.Delete(&model.Category{}, "id = ?", categoryID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *Repository) FindCategoriesByStoreID(storeID uuid.UUID) ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.Where("store_id = ?", storeID).Find(&categories).Error
	return categories, err
}
