package company

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
)

type Service interface {
	GetCompanyByID(ctx context.Context, id string) (*model.Company, error)
	GetAllCompanies(ctx context.Context, search string) ([]*model.Company, error)
	CreateCompany(ctx context.Context, name, address string) (*model.Company, error)
}

type companyService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &companyService{repo: repo}
}

func (s *companyService) GetCompanyByID(ctx context.Context, id string) (*model.Company, error) {
	companyUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, companyUUID)
}

func (s *companyService) GetAllCompanies(ctx context.Context, search string) ([]*model.Company, error) {
	return s.repo.FindAll(ctx, search)
}

func (s *companyService) CreateCompany(ctx context.Context, name, address string) (*model.Company, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("C%04d", r.Intn(10000))

	newCompany := &model.Company{
		Name:    name,
		Address: address,
		Code:    code,
	}

	if err := s.repo.Create(ctx, nil, newCompany); err != nil {
		return nil, err
	}
	return newCompany, nil
}
