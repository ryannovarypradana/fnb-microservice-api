package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	categoryID, err := uuid.Parse(req.GetCategoryId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Category ID format: %v", err)
	}
	storeID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Store ID format: %v", err)
	}

	menu := &model.Menu{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		ImageURL:    req.GetImageUrl(),
		CategoryID:  categoryID,
		StoreID:     storeID,
	}

	if err := h.service.CreateMenu(ctx, menu); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create menu: %v", err)
	}

	return &product.CreateMenuResponse{Menu: modelMenuToProtoMenu(menu)}, nil
}

func (h *ProductGRPCHandler) GetMenuByID(ctx context.Context, req *product.GetMenuByIDRequest) (*product.GetMenuByIDResponse, error) {
	menuID, err := uuid.Parse(req.GetMenuId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Menu ID format: %v", err)
	}

	menu, err := h.service.GetMenuByID(ctx, menuID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Menu not found: %v", err)
	}
	return &product.GetMenuByIDResponse{Menu: modelMenuToProtoMenu(menu)}, nil
}

func (h *ProductGRPCHandler) UpdateMenu(ctx context.Context, req *product.UpdateMenuRequest) (*product.UpdateMenuResponse, error) {
	menuID, err := uuid.Parse(req.GetMenuId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Menu ID format: %v", err)
	}
	categoryID, err := uuid.Parse(req.GetCategoryId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Category ID format: %v", err)
	}

	updatedMenu := &model.Menu{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		ImageURL:    req.GetImageUrl(),
		CategoryID:  categoryID,
	}

	if err := h.service.UpdateMenu(ctx, menuID, updatedMenu); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update menu: %v", err)
	}

	menu, err := h.service.GetMenuByID(ctx, menuID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not retrieve updated menu: %v", err)
	}

	return &product.UpdateMenuResponse{Menu: modelMenuToProtoMenu(menu)}, nil
}

func (h *ProductGRPCHandler) DeleteMenu(ctx context.Context, req *product.DeleteMenuRequest) (*product.DeleteMenuResponse, error) {
	menuID, err := uuid.Parse(req.GetMenuId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Menu ID format: %v", err)
	}
	if err := h.service.DeleteMenu(ctx, menuID); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete menu: %v", err)
	}
	return &product.DeleteMenuResponse{Message: "Menu deleted successfully"}, nil
}

func (h *ProductGRPCHandler) GetMenusByStoreID(ctx context.Context, req *product.GetMenusByStoreIDRequest) (*product.GetMenusByStoreIDResponse, error) {
	storeID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Store ID format: %v", err)
	}

	menus, err := h.service.GetMenusByStoreID(ctx, storeID)
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
	storeID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Store ID format: %v", err)
	}

	category := &model.Category{
		Name:    req.GetName(),
		StoreID: storeID,
	}
	if err := h.service.CreateCategory(ctx, category); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create category: %v", err)
	}
	return &product.CreateCategoryResponse{Category: modelCategoryToProtoCategory(category)}, nil
}

func (h *ProductGRPCHandler) GetCategoryByID(ctx context.Context, req *product.GetCategoryByIDRequest) (*product.GetCategoryByIDResponse, error) {
	categoryID, err := uuid.Parse(req.GetCategoryId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Category ID format: %v", err)
	}

	category, err := h.service.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Category not found: %v", err)
	}
	return &product.GetCategoryByIDResponse{Category: modelCategoryToProtoCategory(category)}, nil
}

func (h *ProductGRPCHandler) UpdateCategory(ctx context.Context, req *product.UpdateCategoryRequest) (*product.UpdateCategoryResponse, error) {
	categoryID, err := uuid.Parse(req.GetCategoryId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Category ID format: %v", err)
	}

	updatedCategory := &model.Category{Name: req.GetName()}
	if err := h.service.UpdateCategory(ctx, categoryID, updatedCategory); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update category: %v", err)
	}

	category, err := h.service.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Could not retrieve updated category: %v", err)
	}

	return &product.UpdateCategoryResponse{Category: modelCategoryToProtoCategory(category)}, nil
}

func (h *ProductGRPCHandler) DeleteCategory(ctx context.Context, req *product.DeleteCategoryRequest) (*product.DeleteCategoryResponse, error) {
	categoryID, err := uuid.Parse(req.GetCategoryId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Category ID format: %v", err)
	}
	if err := h.service.DeleteCategory(ctx, categoryID); err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete category: %v", err)
	}
	return &product.DeleteCategoryResponse{Message: "Category deleted successfully"}, nil
}

func (h *ProductGRPCHandler) GetCategoriesByStoreID(ctx context.Context, req *product.GetCategoriesByStoreIDRequest) (*product.GetCategoriesByStoreIDResponse, error) {
	storeID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Store ID format: %v", err)
	}

	categories, err := h.service.GetCategoriesByStoreID(ctx, storeID)
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
	if menu.Category.ID != uuid.Nil {
		categoryName = menu.Category.Name
	}
	return &product.Menu{
		Id:           menu.ID.String(),
		Name:         menu.Name,
		Description:  menu.Description,
		Price:        menu.Price,
		ImageUrl:     menu.ImageURL,
		CategoryId:   menu.CategoryID.String(),
		StoreId:      menu.StoreID.String(),
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
		Id:        category.ID.String(),
		Name:      category.Name,
		StoreId:   category.StoreID.String(),
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
		UpdatedAt: category.UpdatedAt.Format(time.RFC3339),
	}
}
