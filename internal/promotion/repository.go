package promotion

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

// Repository defines the interface for company data operations.
type Repository interface {
	// Discount methods
	CreateDiscount(discount *model.Discount) error
	FindDiscountByID(id uuid.UUID) (*model.Discount, error)
	FindAllDiscounts() ([]*model.Discount, error)
	UpdateDiscount(discount *model.Discount) error
	DeleteDiscount(id uuid.UUID) error

	// Voucher methods
	CreateVoucher(voucher *model.Voucher) error
	FindVoucherByID(id uuid.UUID) (*model.Voucher, error)
	FindAllVouchers() ([]*model.Voucher, error)
	UpdateVoucher(voucher *model.Voucher) error
	DeleteVoucher(id uuid.UUID) error

	// Bundle methods
	CreateBundle(bundle *model.Bundle) error
	FindBundleByID(id uuid.UUID) (*model.Bundle, error)
	FindAllBundles() ([]*model.Bundle, error)
	UpdateBundle(bundle *model.Bundle) error
	DeleteBundle(id uuid.UUID) error
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

func (r *repository) FindDiscountByID(id uuid.UUID) (*model.Discount, error) {
	var discount model.Discount
	if err := r.db.First(&discount, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("discount not found")
		}
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

func (r *repository) DeleteDiscount(id uuid.UUID) error {
	result := r.db.Delete(&model.Discount{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("discount not found")
	}
	return nil
}

// ========== Voucher Implementation ==========
func (r *repository) CreateVoucher(voucher *model.Voucher) error {
	return r.db.Create(voucher).Error
}

func (r *repository) FindVoucherByID(id uuid.UUID) (*model.Voucher, error) {
	var voucher model.Voucher
	if err := r.db.First(&voucher, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("voucher not found")
		}
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

func (r *repository) DeleteVoucher(id uuid.UUID) error {
	result := r.db.Delete(&model.Voucher{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("voucher not found")
	}
	return nil
}

// ========== Bundle Implementation ==========
func (r *repository) CreateBundle(bundle *model.Bundle) error {
	productIDsJSON, err := json.Marshal(bundle.ProductIDs)
	if err != nil {
		return err
	}
	bundle.Products = string(productIDsJSON)
	return r.db.Create(bundle).Error
}

func (r *repository) FindBundleByID(id uuid.UUID) (*model.Bundle, error) {
	var bundle model.Bundle
	if err := r.db.First(&bundle, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bundle not found")
		}
		return nil, err
	}
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
			return nil, err
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

func (r *repository) DeleteBundle(id uuid.UUID) error {
	result := r.db.Delete(&model.Bundle{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("bundle not found")
	}
	return nil
}
