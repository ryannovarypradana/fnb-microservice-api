package order

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/eventbus"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type Service interface {
	CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*model.Order, error)
	GetOrderByID(ctx context.Context, id string) (*model.Order, error)
	GetAllOrdersByUserID(ctx context.Context, userID string) ([]*model.Order, error)
}

type orderService struct {
	repo           Repository
	productClient  product.ProductServiceClient
	eventPublisher eventbus.Publisher
}

func NewService(repo Repository, productClient product.ProductServiceClient, publisher eventbus.Publisher) Service {
	return &orderService{
		repo:           repo,
		productClient:  productClient,
		eventPublisher: publisher,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*model.Order, error) {
	userUUID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	var orderItems []model.OrderItem
	var totalPrice float64
	for _, item := range req.GetItems() {
		productInfo, err := s.productClient.GetProduct(ctx, &product.GetProductRequest{Id: item.GetProductId()})
		if err != nil {
			return nil, fmt.Errorf("product with id %s not found", item.GetProductId())
		}

		price := productInfo.GetProduct().GetPrice()
		subtotal := price * float64(item.GetQuantity())
		totalPrice += subtotal

		orderItems = append(orderItems, model.OrderItem{
			ProductID: uuid.MustParse(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
			Subtotal:  subtotal,
		})
	}

	newOrder := &model.Order{
		UserID:     userUUID,
		Status:     "PENDING",
		TotalPrice: totalPrice,
	}

	if err := s.repo.Create(ctx, newOrder, orderItems); err != nil {
		return nil, err
	}

	log.Printf("Mempublikasikan event 'order.created' untuk Order ID: %s", newOrder.ID.String())
	eventPayload := map[string]interface{}{
		"order_id":    newOrder.ID.String(),
		"user_id":     newOrder.UserID.String(),
		"total_price": newOrder.TotalPrice,
		"status":      newOrder.Status,
	}

	go func() {
		err := s.eventPublisher.Publish("fnb_events", "order.created", "application/json", eventPayload)
		if err != nil {
			log.Printf("ERROR: Gagal mempublikasikan event order.created untuk Order ID %s: %v", newOrder.ID.String(), err)
		}
	}()

	newOrder.Items = orderItems
	return newOrder, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) (*model.Order, error) {
	orderUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, orderUUID)
}

func (s *orderService) GetAllOrdersByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAllByUserID(ctx, userUUID)
}
