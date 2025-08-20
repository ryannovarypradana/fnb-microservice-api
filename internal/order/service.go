// internal/order/service.go

package order

import (
	"context"
	"errors"
	"fmt"
	"log"

	// Menggunakan tipe klien gRPC yang konkret, bukan interface
	productpb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type IService interface {
	CreateOrder(ctx context.Context, order *model.Order, items []*model.OrderItem) (*model.Order, error)
	GetOrderByID(ctx context.Context, orderID uint) (*model.Order, error)
	GetAllOrdersByUserID(ctx context.Context, userID uint) ([]*model.Order, error)
}

type Service struct {
	repo          IRepository
	productClient productpb.ProductServiceClient // <-- Tipe yang benar
}

// Menggunakan tipe klien gRPC yang konkret di constructor
func NewOrderService(repo IRepository, productClient productpb.ProductServiceClient) IService {
	return &Service{
		repo:          repo,
		productClient: productClient,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *model.Order, items []*model.OrderItem) (*model.Order, error) {
	var totalAmount float64 = 0

	for _, item := range items {
		// Memanggil method gRPC yang benar: GetMenuByID
		grpcRequest := &productpb.GetMenuByIDRequest{MenuId: int64(item.ProductID)}
		productResponse, err := s.productClient.GetMenuByID(ctx, grpcRequest)

		if err != nil {
			log.Printf("Failed to get menu with ID %d: %v", item.ProductID, err)
			return nil, fmt.Errorf("invalid menu item with ID: %d", item.ProductID)
		}

		item.Price = productResponse.Menu.Price
		totalAmount += item.Price * float64(item.Quantity)
	}

	order.TotalAmount = totalAmount
	order.Status = "PENDING" // Set status awal

	if err := s.repo.CreateOrderWithItems(order, items); err != nil {
		return nil, err
	}

	// Mengembalikan order yang sudah dibuat agar handler bisa mendapatkan ID-nya
	return order, nil
}

func (s *Service) GetOrderByID(ctx context.Context, orderID uint) (*model.Order, error) {
	order, err := s.repo.FindOrderByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *Service) GetAllOrdersByUserID(ctx context.Context, userID uint) ([]*model.Order, error) {
	return s.repo.FindOrdersByUserID(userID)
}
