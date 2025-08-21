package promotion

import (
	"context"

	promotionpb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/promotion"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcHandler struct {
	promotionpb.UnimplementedPromotionServiceServer
	service Service
}

func NewGrpcHandler(service Service) *GrpcHandler {
	return &GrpcHandler{service: service}
}

// ========== Discount Handler Implementation ==========
func (h *GrpcHandler) CreateDiscount(ctx context.Context, req *promotionpb.CreateDiscountRequest) (*promotionpb.Discount, error) {
	discount, err := h.service.CreateDiscount(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create discount: %v", err)
	}
	return ModelToProtoDiscount(discount), nil
}

func (h *GrpcHandler) GetDiscount(ctx context.Context, req *promotionpb.GetByIDRequest) (*promotionpb.Discount, error) {
	discount, err := h.service.GetDiscountByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "discount not found: %v", err)
	}
	return ModelToProtoDiscount(discount), nil
}

func (h *GrpcHandler) UpdateDiscount(ctx context.Context, req *promotionpb.UpdateDiscountRequest) (*promotionpb.Discount, error) {
	discount, err := h.service.UpdateDiscount(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update discount: %v", err)
	}
	return ModelToProtoDiscount(discount), nil
}

func (h *GrpcHandler) DeleteDiscount(ctx context.Context, req *promotionpb.GetByIDRequest) (*promotionpb.DeleteResponse, error) {
	if err := h.service.DeleteDiscount(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete discount: %v", err)
	}
	return &promotionpb.DeleteResponse{Message: "Discount deleted successfully"}, nil
}

func (h *GrpcHandler) ListDiscounts(ctx context.Context, _ *emptypb.Empty) (*promotionpb.ListDiscountsResponse, error) {
	discounts, err := h.service.GetAllDiscounts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list discounts: %v", err)
	}
	var pbDiscounts []*promotionpb.Discount
	for _, d := range discounts {
		pbDiscounts = append(pbDiscounts, ModelToProtoDiscount(d))
	}
	return &promotionpb.ListDiscountsResponse{Discounts: pbDiscounts}, nil
}

// ========== Voucher Handler Implementation ==========
func (h *GrpcHandler) CreateVoucher(ctx context.Context, req *promotionpb.CreateVoucherRequest) (*promotionpb.Voucher, error) {
	voucher, err := h.service.CreateVoucher(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create voucher: %v", err)
	}
	return ModelToProtoVoucher(voucher), nil
}

func (h *GrpcHandler) GetVoucher(ctx context.Context, req *promotionpb.GetByIDRequest) (*promotionpb.Voucher, error) {
	voucher, err := h.service.GetVoucherByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "voucher not found: %v", err)
	}
	return ModelToProtoVoucher(voucher), nil
}

func (h *GrpcHandler) UpdateVoucher(ctx context.Context, req *promotionpb.UpdateVoucherRequest) (*promotionpb.Voucher, error) {
	voucher, err := h.service.UpdateVoucher(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update voucher: %v", err)
	}
	return ModelToProtoVoucher(voucher), nil
}

func (h *GrpcHandler) DeleteVoucher(ctx context.Context, req *promotionpb.GetByIDRequest) (*promotionpb.DeleteResponse, error) {
	if err := h.service.DeleteVoucher(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete voucher: %v", err)
	}
	return &promotionpb.DeleteResponse{Message: "Voucher deleted successfully"}, nil
}

func (h *GrpcHandler) ListVouchers(ctx context.Context, _ *emptypb.Empty) (*promotionpb.ListVouchersResponse, error) {
	vouchers, err := h.service.GetAllVouchers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list vouchers: %v", err)
	}
	var pbVouchers []*promotionpb.Voucher
	for _, v := range vouchers {
		pbVouchers = append(pbVouchers, ModelToProtoVoucher(v))
	}
	return &promotionpb.ListVouchersResponse{Vouchers: pbVouchers}, nil
}

// ========== Bundle Handler Implementation ==========
func (h *GrpcHandler) CreateBundle(ctx context.Context, req *promotionpb.CreateBundleRequest) (*promotionpb.Bundle, error) {
	bundle, err := h.service.CreateBundle(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create bundle: %v", err)
	}
	return ModelToProtoBundle(bundle), nil
}

func (h *GrpcHandler) GetBundle(ctx context.Context, req *promotionpb.GetByIDRequest) (*promotionpb.Bundle, error) {
	bundle, err := h.service.GetBundleByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "bundle not found: %v", err)
	}
	return ModelToProtoBundle(bundle), nil
}

func (h *GrpcHandler) UpdateBundle(ctx context.Context, req *promotionpb.UpdateBundleRequest) (*promotionpb.Bundle, error) {
	bundle, err := h.service.UpdateBundle(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update bundle: %v", err)
	}
	return ModelToProtoBundle(bundle), nil
}

func (h *GrpcHandler) DeleteBundle(ctx context.Context, req *promotionpb.GetByIDRequest) (*promotionpb.DeleteResponse, error) {
	if err := h.service.DeleteBundle(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete bundle: %v", err)
	}
	return &promotionpb.DeleteResponse{Message: "Bundle deleted successfully"}, nil
}

func (h *GrpcHandler) ListBundles(ctx context.Context, _ *emptypb.Empty) (*promotionpb.ListBundlesResponse, error) {
	bundles, err := h.service.GetAllBundles(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list bundles: %v", err)
	}
	var pbBundles []*promotionpb.Bundle
	for _, b := range bundles {
		pbBundles = append(pbBundles, ModelToProtoBundle(b))
	}
	return &promotionpb.ListBundlesResponse{Bundles: pbBundles}, nil
}
