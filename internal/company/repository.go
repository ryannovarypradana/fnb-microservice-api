// /internal/company/repository.go
package company

import (
	"context"
	"fmt"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for company data operations.
type Repository interface {
	Create(ctx context.Context, tx *gorm.DB, company *model.Company) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Company, error)
	FindAll(ctx context.Context, search string) ([]*model.Company, error)
	FindByCode(ctx context.Context, code string) (*model.Company, error)
}

type companyRepository struct {
	db *gorm.DB
}

// NewRepository creates a new instance of the company repository.
func NewRepository(db *gorm.DB) Repository {
	return &companyRepository{db: db}
}

// Create creates a new company record in the database.
// It can run within a transaction if tx is not nil.
func (r *companyRepository) Create(ctx context.Context, tx *gorm.DB, company *model.Company) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.WithContext(ctx).Create(company).Error
}

// FindByID retrieves a company by its UUID.
func (r *companyRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	var company model.Company
	if err := r.db.WithContext(ctx).First(&company, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

// FindAll retrieves all companies, with an optional search query on the name.
func (r *companyRepository) FindAll(ctx context.Context, search string) ([]*model.Company, error) {
	var companies []*model.Company
	query := r.db.WithContext(ctx)
	if search != "" {
		query = query.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", search))
	}
	if err := query.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

// FindByCode retrieves a company by its unique code.
func (r *companyRepository) FindByCode(ctx context.Context, code string) (*model.Company, error) {
	var company model.Company
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}
