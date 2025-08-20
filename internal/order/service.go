package order

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/client"
	productpb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type IService interface {
	CreateOrder(ctx context.Context, order *model.Order, items []*model.OrderItem) error
	GetOrderByID(ctx context.Context, orderID uint) (*model.Order, error)
	GetAllOrdersByUserID(ctx context.Context, userID uint) ([]*model.Order, error)
}

type Service struct {
	repo          IRepository
	productClient client.IProductServiceClient
}

func NewOrderService(repo IRepository, productClient client.IProductServiceClient) IService {
	return &Service{
		repo:          repo,
		productClient: productClient,
	}
}

func (s *Service) CreateOrder(ctx context.Context, order *model.Order, items []*model.OrderItem) error {
	var totalAmount float64 = 0

	for _, item := range items {
		grpcRequest := &productpb.GetMenuByIDRequest{MenuId: int64(item.ProductID)}
		productResponse, err := s.productClient.GetProductServiceClient().GetMenuByID(ctx, grpcRequest)

		if err != nil {
			log.Printf("Failed to get menu with ID %d: %v", item.ProductID, err)
			return fmt.Errorf("invalid menu item with ID: %d", item.ProductID)
		}

		item.Price = productResponse.Menu.Price
		totalAmount += item.Price * float64(item.Quantity)
	}

	order.TotalAmount = totalAmount

	if err := s.repo.CreateOrderWithItems(order, items); err != nil {
		return err
	}

	return nil
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
