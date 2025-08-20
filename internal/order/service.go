package order

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	productPB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	storePB "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*model.Order, error)
	GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*model.Order, error)
	CalculateBill(ctx context.Context, req *pb.CalculateBillRequest) (*pb.BillResponse, error)
	CreatePublicOrder(ctx context.Context, req *pb.CreatePublicOrderRequest) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*model.Order, error)
	UpdateOrderItems(ctx context.Context, req *pb.UpdateOrderItemsRequest) (*model.Order, error)
	ConfirmPayment(ctx context.Context, req *pb.ConfirmPaymentRequest) (*model.Order, error)
}

type orderService struct {
	repo          OrderRepository
	db            *gorm.DB
	productClient productPB.ProductServiceClient
	storeClient   storePB.StoreServiceClient
}

func NewOrderService(repo OrderRepository, db *gorm.DB, productClient productPB.ProductServiceClient, storeClient storePB.StoreServiceClient) OrderService {
	return &orderService{repo, db, productClient, storeClient}
}

func (s *orderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*model.Order, error) {
	storeID, err := uuid.Parse(req.StoreId)
	if err != nil {
		return nil, errors.New("invalid store id")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	var orderItems []*model.OrderItem
	var totalAmount float64

	for _, item := range req.Items {
		// Asumsi product.proto menggunakan int64 untuk MenuId
		menuIDInt, err := strconv.ParseInt(item.ProductId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid product id format: %s", item.ProductId)
		}

		product, err := s.productClient.GetMenuByID(ctx, &productPB.GetMenuByIDRequest{MenuId: menuIDInt})
		if err != nil {
			return nil, fmt.Errorf("product %s not found", item.ProductId)
		}

		// Asumsi model.OrderItem menggunakan uuid.UUID untuk ProductID
		menuIDUUID, err := uuid.Parse(strconv.FormatInt(product.Menu.Id, 10))
		if err != nil {
			return nil, fmt.Errorf("cannot parse product id from response: %v", err)
		}

		orderItems = append(orderItems, &model.OrderItem{
			ProductID: menuIDUUID,
			Quantity:  int(item.Quantity),
			Price:     product.Menu.Price,
		})
		totalAmount += product.Menu.Price * float64(item.Quantity)
	}

	newOrder := &model.Order{
		StoreID:     storeID,
		UserID:      &userID,
		TotalAmount: totalAmount,
		Status:      model.StatusPending,
		OrderCode:   fmt.Sprintf("%d", time.Now().UnixNano()+int64(rand.Intn(100))),
		Items:       orderItems,
	}

	if err := s.repo.Create(ctx, newOrder); err != nil {
		return nil, err
	}

	return newOrder, nil
}

func (s *orderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*model.Order, error) {
	return s.repo.FindByID(ctx, req.Id)
}

func (s *orderService) CalculateBill(ctx context.Context, req *pb.CalculateBillRequest) (*pb.BillResponse, error) {
	var subtotal float64
	for _, item := range req.Items {
		menuIDInt, err := strconv.ParseInt(item.ProductId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid product id format: %s", item.ProductId)
		}

		product, err := s.productClient.GetMenuByID(ctx, &productPB.GetMenuByIDRequest{MenuId: menuIDInt})
		if err != nil {
			return nil, fmt.Errorf("product with id %s not found", item.ProductId)
		}
		subtotal += product.Menu.Price * float64(item.Quantity)
	}

	tax := subtotal * 0.11 // Contoh pajak 11%

	return &pb.BillResponse{
		Subtotal:    subtotal,
		Tax:         tax,
		Discount:    0,
		TotalAmount: subtotal + tax,
	}, nil
}

func (s *orderService) CreatePublicOrder(ctx context.Context, req *pb.CreatePublicOrderRequest) (*model.Order, error) {
	// Anda perlu GetStoreByCode di store.proto & service, asumsikan itu ada
	// store, err := s.storeClient.GetStoreByCode(ctx, &storePB.GetStoreByCodeRequest{Code: req.StoreCode})
	// if err != nil {
	// 	return nil, errors.New("store not found")
	// }
	// storeID, _ := uuid.Parse(store.Store.Id)

	storeID := uuid.New() // Placeholder jika GetStoreByCode belum ada

	var orderItems []*model.OrderItem
	var totalAmount float64

	for _, item := range req.Items {
		menuIDInt, err := strconv.ParseInt(item.ProductId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid product id format: %s", item.ProductId)
		}
		product, err := s.productClient.GetMenuByID(ctx, &productPB.GetMenuByIDRequest{MenuId: menuIDInt})
		if err != nil {
			return nil, fmt.Errorf("product %s not found", item.ProductId)
		}

		menuIDUUID, _ := uuid.Parse(strconv.FormatInt(product.Menu.Id, 10))
		orderItems = append(orderItems, &model.OrderItem{
			ProductID: menuIDUUID,
			Quantity:  int(item.Quantity),
			Price:     product.Menu.Price,
		})
		totalAmount += product.Menu.Price * float64(item.Quantity)
	}

	newOrder := &model.Order{
		StoreID:      storeID,
		TotalAmount:  totalAmount,
		Status:       model.StatusPending,
		OrderCode:    fmt.Sprintf("%d", time.Now().UnixNano()+int64(rand.Intn(100))),
		CustomerName: req.CustomerName,
		TableNumber:  req.TableNumber,
		Items:        orderItems,
	}

	if err := s.repo.Create(ctx, newOrder); err != nil {
		return nil, err
	}

	return newOrder, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*model.Order, error) {
	order, err := s.repo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	order.Status = model.OrderStatus(req.Status)
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, req.Id)
}

func (s *orderService) UpdateOrderItems(ctx context.Context, req *pb.UpdateOrderItemsRequest) (*model.Order, error) {
	var updatedOrder *model.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		order, err := s.repo.FindByID(ctx, req.Id)
		if err != nil {
			return err
		}

		if order.Status != model.StatusPending {
			return errors.New("order can only be updated if status is pending")
		}

		if err := tx.Where("order_id = ?", order.ID).Delete(&model.OrderItem{}).Error; err != nil {
			return err
		}

		var newItems []*model.OrderItem
		var newTotal float64
		for _, item := range req.Items {
			menuIDInt, err := strconv.ParseInt(item.ProductId, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid product id format: %s", item.ProductId)
			}
			product, err := s.productClient.GetMenuByID(ctx, &productPB.GetMenuByIDRequest{MenuId: menuIDInt})
			if err != nil {
				return fmt.Errorf("product %s not found", item.ProductId)
			}
			menuIDUUID, _ := uuid.Parse(strconv.FormatInt(product.Menu.Id, 10))
			newItem := &model.OrderItem{
				OrderID:   order.ID,
				ProductID: menuIDUUID,
				Quantity:  int(item.Quantity),
				Price:     product.Menu.Price,
			}
			newItems = append(newItems, newItem)
			newTotal += product.Menu.Price * float64(item.Quantity)
		}

		if err := tx.Create(&newItems).Error; err != nil {
			return err
		}

		order.TotalAmount = newTotal
		if err := s.repo.Update(ctx, order); err != nil {
			return err
		}

		updatedOrder = order
		updatedOrder.Items = newItems
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, updatedOrder.ID.String())
}

func (s *orderService) ConfirmPayment(ctx context.Context, req *pb.ConfirmPaymentRequest) (*model.Order, error) {
	order, err := s.repo.FindByCode(ctx, req.OrderCode)
	if err != nil {
		return nil, err
	}
	if order.Status != model.StatusPending {
		return nil, errors.New("order is not pending")
	}
	order.Status = model.StatusPaid
	if err := s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	return s.repo.FindByCode(ctx, req.OrderCode)
}
