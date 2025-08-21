package promotion

import (
	"encoding/json"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	// Discount methods
	CreateDiscount(discount *model.Discount) error
	FindDiscountByID(id uint64) (*model.Discount, error)
	FindAllDiscounts() ([]*model.Discount, error)
	UpdateDiscount(discount *model.Discount) error
	DeleteDiscount(id uint64) error

	// Voucher methods
	CreateVoucher(voucher *model.Voucher) error
	FindVoucherByID(id uint64) (*model.Voucher, error)
	FindAllVouchers() ([]*model.Voucher, error)
	UpdateVoucher(voucher *model.Voucher) error
	DeleteVoucher(id uint64) error

	// Bundle methods
	CreateBundle(bundle *model.Bundle) error
	FindBundleByID(id uint64) (*model.Bundle, error)
	FindAllBundles() ([]*model.Bundle, error)
	UpdateBundle(bundle *model.Bundle) error
	DeleteBundle(id uint64) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

// ========== Discount Implementation ==========
func (r *repository) CreateDiscount(discount *model.Discount) error {
	return r.db.Create(discount).Error
}

func (r *repository) FindDiscountByID(id uint64) (*model.Discount, error) {
	var discount model.Discount
	if err := r.db.First(&discount, id).Error; err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *repository) FindAllDiscounts() ([]*model.Discount, error) {
	var discounts []*model.Discount
	if err := r.db.Find(&discounts).Error; err != nil {
		return nil, err
	}
	return discounts, nil
}

func (r *repository) UpdateDiscount(discount *model.Discount) error {
	return r.db.Save(discount).Error
}

func (r *repository) DeleteDiscount(id uint64) error {
	return r.db.Delete(&model.Discount{}, id).Error
}

// ========== Voucher Implementation ==========
func (r *repository) CreateVoucher(voucher *model.Voucher) error {
	return r.db.Create(voucher).Error
}

func (r *repository) FindVoucherByID(id uint64) (*model.Voucher, error) {
	var voucher model.Voucher
	if err := r.db.First(&voucher, id).Error; err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *repository) FindAllVouchers() ([]*model.Voucher, error) {
	var vouchers []*model.Voucher
	if err := r.db.Find(&vouchers).Error; err != nil {
		return nil, err
	}
	return vouchers, nil
}

func (r *repository) UpdateVoucher(voucher *model.Voucher) error {
	return r.db.Save(voucher).Error
}

func (r *repository) DeleteVoucher(id uint64) error {
	return r.db.Delete(&model.Voucher{}, id).Error
}

// ========== Bundle Implementation ==========
func (r *repository) CreateBundle(bundle *model.Bundle) error {
	// Marshal ProductIDs to JSON string before saving
	productIDsJSON, err := json.Marshal(bundle.ProductIDs)
	if err != nil {
		return err
	}
	bundle.Products = string(productIDsJSON)
	return r.db.Create(bundle).Error
}

func (r *repository) FindBundleByID(id uint64) (*model.Bundle, error) {
	var bundle model.Bundle
	if err := r.db.First(&bundle, id).Error; err != nil {
		return nil, err
	}
	// Unmarshal JSON string to ProductIDs after fetching
	if err := json.Unmarshal([]byte(bundle.Products), &bundle.ProductIDs); err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (r *repository) FindAllBundles() ([]*model.Bundle, error) {
	var bundles []*model.Bundle
	if err := r.db.Find(&bundles).Error; err != nil {
		return nil, err
	}
	for _, b := range bundles {
		if err := json.Unmarshal([]byte(b.Products), &b.ProductIDs); err != nil {
			return nil, err // Or handle error more gracefully
		}
	}
	return bundles, nil
}

func (r *repository) UpdateBundle(bundle *model.Bundle) error {
	productIDsJSON, err := json.Marshal(bundle.ProductIDs)
	if err != nil {
		return err
	}
	bundle.Products = string(productIDsJSON)
	return r.db.Save(bundle).Error
}

func (r *repository) DeleteBundle(id uint64) error {
	return r.db.Delete(&model.Bundle{}, id).Error
}
