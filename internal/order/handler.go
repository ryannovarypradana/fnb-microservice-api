// internal/order/handler.go

package order

import (
	"context"

	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGRPCHandler struct {
	pb.UnimplementedOrderServiceServer
	service IService
}

// INILAH FUNGSI YANG HARUS ADA
func NewOrderGRPCHandler(service IService) *OrderGRPCHandler {
	return &OrderGRPCHandler{service: service}
}

func (h *OrderGRPCHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	orderModel := &model.Order{
		UserID:  uint(req.GetUserId()),
		StoreID: uint(req.GetStoreId()),
	}

	var items []*model.OrderItem
	for _, itemReq := range req.GetItems() {
		items = append(items, &model.OrderItem{
			ProductID: uint(itemReq.GetProductId()),
			Quantity:  int(itemReq.GetQuantity()),
		})
	}

	// Memanggil service dan menampung order yang sudah dibuat
	createdOrder, err := h.service.CreateOrder(ctx, orderModel, items)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create order: %v", err)
	}

	return &pb.CreateOrderResponse{
		Id:          int64(createdOrder.ID),
		Status:      createdOrder.Status,
		TotalAmount: createdOrder.TotalAmount,
	}, nil
}

func (h *OrderGRPCHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := h.service.GetOrderByID(ctx, uint(req.GetId()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Order not found: %v", err)
	}

	return &pb.GetOrderResponse{
		Id:          int64(order.ID),
		UserId:      int64(order.UserID),
		StoreId:     int64(order.StoreID),
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
	}, nil
}

func (h *OrderGRPCHandler) GetAllOrders(ctx context.Context, req *pb.GetAllOrdersRequest) (*pb.GetAllOrdersResponse, error) {
	orders, err := h.service.GetAllOrdersByUserID(ctx, uint(req.GetUserId()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get orders: %v", err)
	}

	var orderResponses []*pb.GetOrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, &pb.GetOrderResponse{
			Id:          int64(order.ID),
			UserId:      int64(order.UserID),
			StoreId:     int64(order.StoreID),
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
		})
	}

	return &pb.GetAllOrdersResponse{Orders: orderResponses}, nil
}
