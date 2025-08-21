package promotion

import (
	"context"

	promotionpb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/promotion"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	// Discount methods
	CreateDiscount(ctx context.Context, req *promotionpb.CreateDiscountRequest) (*model.Discount, error)
	GetDiscountByID(ctx context.Context, id uint64) (*model.Discount, error)
	UpdateDiscount(ctx context.Context, req *promotionpb.UpdateDiscountRequest) (*model.Discount, error)
	DeleteDiscount(ctx context.Context, id uint64) error
	GetAllDiscounts(ctx context.Context) ([]*model.Discount, error)

	// Voucher methods
	CreateVoucher(ctx context.Context, req *promotionpb.CreateVoucherRequest) (*model.Voucher, error)
	GetVoucherByID(ctx context.Context, id uint64) (*model.Voucher, error)
	UpdateVoucher(ctx context.Context, req *promotionpb.UpdateVoucherRequest) (*model.Voucher, error)
	DeleteVoucher(ctx context.Context, id uint64) error
	GetAllVouchers(ctx context.Context) ([]*model.Voucher, error)

	// Bundle methods
	CreateBundle(ctx context.Context, req *promotionpb.CreateBundleRequest) (*model.Bundle, error)
	GetBundleByID(ctx context.Context, id uint64) (*model.Bundle, error)
	UpdateBundle(ctx context.Context, req *promotionpb.UpdateBundleRequest) (*model.Bundle, error)
	DeleteBundle(ctx context.Context, id uint64) error
	GetAllBundles(ctx context.Context) ([]*model.Bundle, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

// ========== Discount Service Implementation ==========
func (s *service) CreateDiscount(ctx context.Context, req *promotionpb.CreateDiscountRequest) (*model.Discount, error) {
	discount := &model.Discount{
		Name:      req.Name,
		Type:      req.Type,
		Value:     req.Value,
		StartDate: req.StartDate.AsTime(),
		EndDate:   req.EndDate.AsTime(),
	}
	if err := s.repository.CreateDiscount(discount); err != nil {
		return nil, err
	}
	return discount, nil
}

func (s *service) GetDiscountByID(ctx context.Context, id uint64) (*model.Discount, error) {
	return s.repository.FindDiscountByID(id)
}

func (s *service) UpdateDiscount(ctx context.Context, req *promotionpb.UpdateDiscountRequest) (*model.Discount, error) {
	discount, err := s.repository.FindDiscountByID(req.Id)
	if err != nil {
		return nil, err
	}
	discount.Name = req.Name
	discount.Type = req.Type
	discount.Value = req.Value
	discount.StartDate = req.StartDate.AsTime()
	discount.EndDate = req.EndDate.AsTime()

	if err := s.repository.UpdateDiscount(discount); err != nil {
		return nil, err
	}
	return discount, nil
}

func (s *service) DeleteDiscount(ctx context.Context, id uint64) error {
	return s.repository.DeleteDiscount(id)
}

func (s *service) GetAllDiscounts(ctx context.Context) ([]*model.Discount, error) {
	return s.repository.FindAllDiscounts()
}

// ========== Voucher Service Implementation ==========
func (s *service) CreateVoucher(ctx context.Context, req *promotionpb.CreateVoucherRequest) (*model.Voucher, error) {
	voucher := &model.Voucher{
		Code:      req.Code,
		Type:      req.Type,
		Value:     req.Value,
		Quota:     int(req.Quota),
		StartDate: req.StartDate.AsTime(),
		EndDate:   req.EndDate.AsTime(),
	}
	if err := s.repository.CreateVoucher(voucher); err != nil {
		return nil, err
	}
	return voucher, nil
}

func (s *service) GetVoucherByID(ctx context.Context, id uint64) (*model.Voucher, error) {
	return s.repository.FindVoucherByID(id)
}

func (s *service) UpdateVoucher(ctx context.Context, req *promotionpb.UpdateVoucherRequest) (*model.Voucher, error) {
	voucher, err := s.repository.FindVoucherByID(req.Id)
	if err != nil {
		return nil, err
	}
	voucher.Code = req.Code
	voucher.Type = req.Type
	voucher.Value = req.Value
	voucher.Quota = int(req.Quota)
	voucher.StartDate = req.StartDate.AsTime()
	voucher.EndDate = req.EndDate.AsTime()

	if err := s.repository.UpdateVoucher(voucher); err != nil {
		return nil, err
	}
	return voucher, nil
}

func (s *service) DeleteVoucher(ctx context.Context, id uint64) error {
	return s.repository.DeleteVoucher(id)
}

func (s *service) GetAllVouchers(ctx context.Context) ([]*model.Voucher, error) {
	return s.repository.FindAllVouchers()
}

// ========== Bundle Service Implementation ==========
func (s *service) CreateBundle(ctx context.Context, req *promotionpb.CreateBundleRequest) (*model.Bundle, error) {
	productIDs := make([]uint, len(req.ProductIds))
	for i, id := range req.ProductIds {
		productIDs[i] = uint(id)
	}

	bundle := &model.Bundle{
		Name:       req.Name,
		Price:      req.Price,
		ProductIDs: productIDs,
	}

	if err := s.repository.CreateBundle(bundle); err != nil {
		return nil, err
	}
	return bundle, nil
}

func (s *service) GetBundleByID(ctx context.Context, id uint64) (*model.Bundle, error) {
	return s.repository.FindBundleByID(id)
}

func (s *service) UpdateBundle(ctx context.Context, req *promotionpb.UpdateBundleRequest) (*model.Bundle, error) {
	bundle, err := s.repository.FindBundleByID(req.Id)
	if err != nil {
		return nil, err
	}

	productIDs := make([]uint, len(req.ProductIds))
	for i, id := range req.ProductIds {
		productIDs[i] = uint(id)
	}

	bundle.Name = req.Name
	bundle.Price = req.Price
	bundle.ProductIDs = productIDs

	if err := s.repository.UpdateBundle(bundle); err != nil {
		return nil, err
	}
	return bundle, nil
}

func (s *service) DeleteBundle(ctx context.Context, id uint64) error {
	return s.repository.DeleteBundle(id)
}

func (s *service) GetAllBundles(ctx context.Context) ([]*model.Bundle, error) {
	return s.repository.FindAllBundles()
}

// ========== Helper Functions to convert model to proto ==========

func ModelToProtoDiscount(d *model.Discount) *promotionpb.Discount {
	return &promotionpb.Discount{
		Id:        uint64(d.ID),
		Name:      d.Name,
		Type:      d.Type,
		Value:     d.Value,
		StartDate: timestamppb.New(d.StartDate),
		EndDate:   timestamppb.New(d.EndDate),
	}
}

func ModelToProtoVoucher(v *model.Voucher) *promotionpb.Voucher {
	return &promotionpb.Voucher{
		Id:        uint64(v.ID),
		Code:      v.Code,
		Type:      v.Type,
		Value:     v.Value,
		Quota:     int32(v.Quota),
		StartDate: timestamppb.New(v.StartDate),
		EndDate:   timestamppb.New(v.EndDate),
	}
}

func ModelToProtoBundle(b *model.Bundle) *promotionpb.Bundle {
	pids := make([]uint64, len(b.ProductIDs))
	for i, id := range b.ProductIDs {
		pids[i] = uint64(id)
	}
	return &promotionpb.Bundle{
		Id:         uint64(b.ID),
		Name:       b.Name,
		Price:      b.Price,
		ProductIds: pids,
	}
}
