package order

import (
	"context"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type GRPCHandler struct {
	order.UnimplementedOrderServiceServer
	service Service
}

func NewGRPCHandler(s Service) *GRPCHandler {
	return &GRPCHandler{service: s}
}

func (h *GRPCHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	o, err := h.service.CreateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResponse{Order: modelToProto(o)}, nil
}

func (h *GRPCHandler) GetOrderById(ctx context.Context, req *order.GetOrderByIdRequest) (*order.GetOrderByIdResponse, error) {
	o, err := h.service.GetOrderByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &order.GetOrderByIdResponse{Order: modelToProto(o)}, nil
}

func (h *GRPCHandler) GetAllOrders(ctx context.Context, req *order.GetAllOrdersRequest) (*order.GetAllOrdersResponse, error) {
	orders, err := h.service.GetAllOrdersByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	var orderMessages []*order.Order
	for _, o := range orders {
		orderMessages = append(orderMessages, modelToProto(o))
	}
	return &order.GetAllOrdersResponse{Orders: orderMessages}, nil
}

// Helper function to convert internal model to proto message
func modelToProto(o *model.Order) *order.Order {
	var items []*order.OrderItem
	for _, item := range o.Items {
		items = append(items, &order.OrderItem{
			ProductId: item.ProductID.String(),
			Quantity:  int32(item.Quantity),
		})
	}
	return &order.Order{
		Id:         o.ID.String(),
		UserId:     o.UserID.String(),
		Status:     o.Status,
		TotalPrice: o.TotalPrice,
		Items:      items,
	}
}
