package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type IService interface {
	// Menu services
	CreateMenu(ctx context.Context, menu *model.Menu) error
	GetMenuByID(ctx context.Context, menuID uuid.UUID) (*model.Menu, error)
	UpdateMenu(ctx context.Context, menuID uuid.UUID, updatedMenu *model.Menu) error
	DeleteMenu(ctx context.Context, menuID uuid.UUID) error
	GetMenusByStoreID(ctx context.Context, storeID uuid.UUID) ([]*model.Menu, error)

	// Category services
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*model.Category, error)
	UpdateCategory(ctx context.Context, categoryID uuid.UUID, updatedCategory *model.Category) error
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error
	GetCategoriesByStoreID(ctx context.Context, storeID uuid.UUID) ([]*model.Category, error)
}

type Service struct {
	repo IRepository
}

func NewProductService(repo IRepository) IService {
	return &Service{repo: repo}
}

// --- Menu Service Implementations ---

func (s *Service) CreateMenu(ctx context.Context, menu *model.Menu) error {
	return s.repo.CreateMenu(menu)
}

func (s *Service) GetMenuByID(ctx context.Context, menuID uuid.UUID) (*model.Menu, error) {
	return s.repo.GetMenuByID(menuID)
}

func (s *Service) UpdateMenu(ctx context.Context, menuID uuid.UUID, updatedMenu *model.Menu) error {
	menu, err := s.repo.GetMenuByID(menuID)
	if err != nil {
		return err
	}
	menu.Name = updatedMenu.Name
	menu.Description = updatedMenu.Description
	menu.Price = updatedMenu.Price
	menu.ImageURL = updatedMenu.ImageURL
	menu.CategoryID = updatedMenu.CategoryID

	return s.repo.UpdateMenu(menu)
}

func (s *Service) DeleteMenu(ctx context.Context, menuID uuid.UUID) error {
	return s.repo.DeleteMenu(menuID)
}

func (s *Service) GetMenusByStoreID(ctx context.Context, storeID uuid.UUID) ([]*model.Menu, error) {
	return s.repo.FindMenusByStoreID(storeID)
}

// --- Category Service Implementations ---

func (s *Service) CreateCategory(ctx context.Context, category *model.Category) error {
	return s.repo.CreateCategory(category)
}

func (s *Service) GetCategoryByID(ctx context.Context, categoryID uuid.UUID) (*model.Category, error) {
	return s.repo.GetCategoryByID(categoryID)
}

func (s *Service) UpdateCategory(ctx context.Context, categoryID uuid.UUID, updatedCategory *model.Category) error {
	category, err := s.repo.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}
	category.Name = updatedCategory.Name

	return s.repo.UpdateCategory(category)
}

func (s *Service) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return s.repo.DeleteCategory(categoryID)
}

func (s *Service) GetCategoriesByStoreID(ctx context.Context, storeID uuid.UUID) ([]*model.Category, error) {
	return s.repo.FindCategoriesByStoreID(storeID)
}
