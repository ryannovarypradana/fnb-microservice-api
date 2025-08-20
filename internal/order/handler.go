package order

import (
	"context"

	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGRPCHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
}

func NewOrderGRPCHandler(grpcServer *grpc.Server, service OrderService) {
	handler := &OrderGRPCHandler{service: service}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func mapOrderToResponse(order *model.Order) *pb.OrderResponse {
	res := &pb.OrderResponse{
		Id:           order.ID.String(),
		StoreId:      order.StoreID.String(),
		TotalAmount:  order.TotalAmount,
		Status:       string(order.Status),
		OrderCode:    order.OrderCode,
		CustomerName: order.CustomerName,
		TableNumber:  order.TableNumber,
	}
	if order.UserID != nil {
		userID := order.UserID.String()
		res.UserId = userID
	}
	for _, item := range order.Items {
		res.Items = append(res.Items, &pb.OrderItemResponse{
			ProductId:   item.ProductID.String(),
			ProductName: item.Menu.Name, // Menggunakan relasi Menu dari GORM
			Price:       item.Price,
			Quantity:    int32(item.Quantity),
		})
	}
	return res
}

func (h *OrderGRPCHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	order, err := h.service.CreateOrder(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}
	return mapOrderToResponse(order), nil
}

func (h *OrderGRPCHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	order, err := h.service.GetOrder(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}
	return mapOrderToResponse(order), nil
}

func (h *OrderGRPCHandler) GetAllOrders(ctx context.Context, req *pb.GetAllOrdersRequest) (*pb.GetAllOrdersResponse, error) {
	// Logika untuk GetAllOrders perlu diimplementasikan di service dan repository
	return nil, status.Errorf(codes.Unimplemented, "method GetAllOrders not implemented")
}

func (h *OrderGRPCHandler) CalculateBill(ctx context.Context, req *pb.CalculateBillRequest) (*pb.BillResponse, error) {
	bill, err := h.service.CalculateBill(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to calculate bill: %v", err)
	}
	return bill, nil
}

func (h *OrderGRPCHandler) CreatePublicOrder(ctx context.Context, req *pb.CreatePublicOrderRequest) (*pb.OrderResponse, error) {
	order, err := h.service.CreatePublicOrder(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create public order: %v", err)
	}
	return mapOrderToResponse(order), nil
}

func (h *OrderGRPCHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.OrderResponse, error) {
	order, err := h.service.UpdateOrderStatus(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update status: %v", err)
	}
	return mapOrderToResponse(order), nil
}

func (h *OrderGRPCHandler) UpdateOrderItems(ctx context.Context, req *pb.UpdateOrderItemsRequest) (*pb.OrderResponse, error) {
	order, err := h.service.UpdateOrderItems(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update items: %v", err)
	}
	return mapOrderToResponse(order), nil
}

func (h *OrderGRPCHandler) ConfirmPayment(ctx context.Context, req *pb.ConfirmPaymentRequest) (*pb.OrderResponse, error) {
	order, err := h.service.ConfirmPayment(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to confirm payment: %v", err)
	}
	return mapOrderToResponse(order), nil
}
