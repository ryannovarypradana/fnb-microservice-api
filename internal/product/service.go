package product

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store" // <-- Import proto store
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type Service interface {
	CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*model.Product, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	GetAllProducts(ctx context.Context, req *product.GetAllProductsRequest) ([]*model.Product, error)
	UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(ctx context.Context, id string) error
}

type productService struct {
	repo        Repository
	storeClient store.StoreServiceClient // <-- Tambahkan field untuk store client
}

// Perbarui constructor untuk menerima storeClient
func NewService(repo Repository, storeClient store.StoreServiceClient) Service {
	return &productService{
		repo:        repo,
		storeClient: storeClient,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (*model.Product, error) {
	storeUUID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		return nil, errors.New("invalid store id format")
	}

	// === PANGGILAN GPRC UNTUK VALIDASI TOKO ===
	_, err = s.storeClient.GetStore(ctx, &store.GetStoreRequest{Id: req.GetStoreId()})
	if err != nil {
		return nil, errors.New("toko yang dituju tidak valid atau tidak ditemukan")
	}
	// =======================================

	newProduct := &model.Product{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		StoreID:     storeUUID,
	}

	if err := s.repo.Create(ctx, newProduct); err != nil {
		return nil, err
	}
	return newProduct, nil
}

// (Sisa fungsi service lainnya tidak berubah)

func (s *productService) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	productUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, productUUID)
}

func (s *productService) GetAllProducts(ctx context.Context, req *product.GetAllProductsRequest) ([]*model.Product, error) {
	return s.repo.FindAll(ctx, req.GetSearch(), req.GetStoreId())
}

func (s *productService) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*model.Product, error) {
	productUUID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, err
	}

	productToUpdate, err := s.repo.FindByID(ctx, productUUID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	productToUpdate.Name = req.GetName()
	productToUpdate.Description = req.GetDescription()
	productToUpdate.Price = req.GetPrice()

	if err := s.repo.Update(ctx, productToUpdate); err != nil {
		return nil, err
	}
	return productToUpdate, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	productUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, productUUID)
}
