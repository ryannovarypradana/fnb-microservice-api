package product

import (
	"context"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
)

type GRPCHandler struct {
	product.UnimplementedProductServiceServer
	service Service
}

func NewGRPCHandler(s Service) *GRPCHandler {
	return &GRPCHandler{service: s}
}

func (h *GRPCHandler) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	p, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return &product.CreateProductResponse{Product: &product.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		StoreId:     p.StoreID.String(),
	}}, nil
}

func (h *GRPCHandler) GetProduct(ctx context.Context, req *product.GetProductRequest) (*product.GetProductResponse, error) {
	p, err := h.service.GetProductByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &product.GetProductResponse{Product: &product.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		StoreId:     p.StoreID.String(),
	}}, nil
}

func (h *GRPCHandler) GetAllProducts(ctx context.Context, req *product.GetAllProductsRequest) (*product.GetAllProductsResponse, error) {
	products, err := h.service.GetAllProducts(ctx, req)
	if err != nil {
		return nil, err
	}
	var productMessages []*product.Product
	for _, p := range products {
		productMessages = append(productMessages, &product.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			StoreId:     p.StoreID.String(),
		})
	}
	return &product.GetAllProductsResponse{Products: productMessages}, nil
}

func (h *GRPCHandler) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	p, err := h.service.UpdateProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return &product.UpdateProductResponse{Product: &product.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		StoreId:     p.StoreID.String(),
	}}, nil
}

func (h *GRPCHandler) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	if err := h.service.DeleteProduct(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &product.DeleteProductResponse{Message: "Product deleted successfully"}, nil
}
