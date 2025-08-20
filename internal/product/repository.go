package product

import (
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type IRepository interface {
	// Menu operations
	CreateMenu(menu *model.Menu) error
	GetMenuByID(menuID uint) (*model.Menu, error)
	UpdateMenu(menu *model.Menu) error
	DeleteMenu(menuID uint) error
	FindMenusByStoreID(storeID uint) ([]*model.Menu, error)

	// Category operations
	CreateCategory(category *model.Category) error
	GetCategoryByID(categoryID uint) (*model.Category, error)
	UpdateCategory(category *model.Category) error
	DeleteCategory(categoryID uint) error
	FindCategoriesByStoreID(storeID uint) ([]*model.Category, error)
}

type Repository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) IRepository {
	return &Repository{db: db}
}

// --- Menu Implementations ---

func (r *Repository) CreateMenu(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *Repository) GetMenuByID(menuID uint) (*model.Menu, error) {
	var menu model.Menu
	if err := r.db.Preload("Category").First(&menu, menuID).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *Repository) UpdateMenu(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

func (r *Repository) DeleteMenu(menuID uint) error {
	return r.db.Delete(&model.Menu{}, menuID).Error
}

func (r *Repository) FindMenusByStoreID(storeID uint) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := r.db.Preload("Category").Where("store_id = ?", storeID).Find(&menus).Error
	return menus, err
}

// --- Category Implementations ---

func (r *Repository) CreateCategory(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *Repository) GetCategoryByID(categoryID uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.First(&category, categoryID).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) UpdateCategory(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *Repository) DeleteCategory(categoryID uint) error {
	// Pastikan tidak ada menu yang terhubung sebelum menghapus
	var count int64
	r.db.Model(&model.Menu{}).Where("category_id = ?", categoryID).Count(&count)
	if count > 0 {
		return gorm.ErrForeignKeyViolated
	}
	return r.db.Delete(&model.Category{}, categoryID).Error
}

func (r *Repository) FindCategoriesByStoreID(storeID uint) ([]*model.Category, error) {
	var categories []*model.Category
	err := r.db.Where("store_id = ?", storeID).Find(&categories).Error
	return categories, err
}
