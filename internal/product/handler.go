package product

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// Blok import yang divalidasi
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type ProductGRPCHandler struct {
	product.UnimplementedProductServiceServer
	service IService
}

func NewProductGRPCHandler(service IService) *ProductGRPCHandler {
	return &ProductGRPCHandler{service: service}
}

// --- Menu RPC Handlers ---

func (h *ProductGRPCHandler) CreateMenu(ctx context.Context, req *product.CreateMenuRequest) (*product.CreateMenuResponse, error) {
	menu := &model.Menu{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageUrl,
		CategoryID:  uint(req.CategoryId),
		StoreID:     uint(req.StoreId),
	}

	if err := h.service.CreateMenu(ctx, menu); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create menu: %v", err)
	}

	return &product.CreateMenuResponse{Menu: modelMenuToProtoMenu(menu)}, nil
}

func (h *ProductGRPCHandler) GetMenuByID(ctx context.Context, req *product.GetMenuByIDRequest) (*product.GetMenuByIDResponse, error) {
	menu, err := h.service.GetMenuByID(ctx, uint(req.MenuId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Menu not found: %v", err)
	}
	return &product.GetMenuByIDResponse{Menu: modelMenuToProtoMenu(menu)}, nil
}

func (h *ProductGRPCHandler) UpdateMenu(ctx context.Context, req *product.UpdateMenuRequest) (*product.UpdateMenuResponse, error) {
	updatedMenu := &model.Menu{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageUrl,
		CategoryID:  uint(req.CategoryId),
	}

	if err := h.service.UpdateMenu(ctx, uint(req.MenuId), updatedMenu); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update menu: %v", err)
	}

	menu, err := h.service.GetMenuByID(ctx, uint(req.MenuId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not retrieve updated menu: %v", err)
	}

	return &product.UpdateMenuResponse{Menu: modelMenuToProtoMenu(menu)}, nil
}

func (h *ProductGRPCHandler) DeleteMenu(ctx context.Context, req *product.DeleteMenuRequest) (*product.DeleteMenuResponse, error) {
	if err := h.service.DeleteMenu(ctx, uint(req.MenuId)); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete menu: %v", err)
	}
	return &product.DeleteMenuResponse{Message: "Menu deleted successfully"}, nil
}

func (h *ProductGRPCHandler) GetMenusByStoreID(ctx context.Context, req *product.GetMenusByStoreIDRequest) (*product.GetMenusByStoreIDResponse, error) {
	menus, err := h.service.GetMenusByStoreID(ctx, uint(req.StoreId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get menus: %v", err)
	}

	var protoMenus []*product.Menu
	for _, menu := range menus {
		protoMenus = append(protoMenus, modelMenuToProtoMenu(menu))
	}

	return &product.GetMenusByStoreIDResponse{Menus: protoMenus}, nil
}

// --- Category RPC Handlers ---

func (h *ProductGRPCHandler) CreateCategory(ctx context.Context, req *product.CreateCategoryRequest) (*product.CreateCategoryResponse, error) {
	category := &model.Category{
		Name:    req.Name,
		StoreID: uint(req.StoreId),
	}
	if err := h.service.CreateCategory(ctx, category); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create category: %v", err)
	}
	return &product.CreateCategoryResponse{Category: modelCategoryToProtoCategory(category)}, nil
}

func (h *ProductGRPCHandler) GetCategoryByID(ctx context.Context, req *product.GetCategoryByIDRequest) (*product.GetCategoryByIDResponse, error) {
	category, err := h.service.GetCategoryByID(ctx, uint(req.CategoryId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Category not found: %v", err)
	}
	return &product.GetCategoryByIDResponse{Category: modelCategoryToProtoCategory(category)}, nil
}

func (h *ProductGRPCHandler) UpdateCategory(ctx context.Context, req *product.UpdateCategoryRequest) (*product.UpdateCategoryResponse, error) {
	updatedCategory := &model.Category{Name: req.Name}
	if err := h.service.UpdateCategory(ctx, uint(req.CategoryId), updatedCategory); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update category: %v", err)
	}

	category, err := h.service.GetCategoryByID(ctx, uint(req.CategoryId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not retrieve updated category: %v", err)
	}

	return &product.UpdateCategoryResponse{Category: modelCategoryToProtoCategory(category)}, nil
}

func (h *ProductGRPCHandler) DeleteCategory(ctx context.Context, req *product.DeleteCategoryRequest) (*product.DeleteCategoryResponse, error) {
	if err := h.service.DeleteCategory(ctx, uint(req.CategoryId)); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete category: %v", err)
	}
	return &product.DeleteCategoryResponse{Message: "Category deleted successfully"}, nil
}

func (h *ProductGRPCHandler) GetCategoriesByStoreID(ctx context.Context, req *product.GetCategoriesByStoreIDRequest) (*product.GetCategoriesByStoreIDResponse, error) {
	categories, err := h.service.GetCategoriesByStoreID(ctx, uint(req.StoreId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get categories: %v", err)
	}

	var protoCategories []*product.Category
	for _, cat := range categories {
		protoCategories = append(protoCategories, modelCategoryToProtoCategory(cat))
	}

	return &product.GetCategoriesByStoreIDResponse{Categories: protoCategories}, nil
}

// --- Helper Functions for Conversion ---
func modelMenuToProtoMenu(menu *model.Menu) *product.Menu {
	if menu == nil {
		return nil
	}
	categoryName := ""
	if menu.Category.ID != 0 {
		categoryName = menu.Category.Name
	}
	return &product.Menu{
		Id:           int64(menu.ID),
		Name:         menu.Name,
		Description:  menu.Description,
		Price:        menu.Price,
		ImageUrl:     menu.ImageURL,
		CategoryId:   int64(menu.CategoryID),
		StoreId:      int64(menu.StoreID),
		CategoryName: categoryName,
		CreatedAt:    menu.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    menu.UpdatedAt.Format(time.RFC3339),
	}
}

func modelCategoryToProtoCategory(category *model.Category) *product.Category {
	if category == nil {
		return nil
	}
	return &product.Category{
		Id:        int64(category.ID),
		Name:      category.Name,
		StoreId:   int64(category.StoreID),
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
		UpdatedAt: category.UpdatedAt.Format(time.RFC3339),
	}
}
