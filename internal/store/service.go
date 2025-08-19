package store

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

// FIX: Tambahkan semua metode yang akan digunakan oleh handler ke dalam interface ini.
type Service interface {
	CreateStore(ctx context.Context, companyID, name, address string) (*model.Store, error)
	GetStoreByID(ctx context.Context, id string) (*model.Store, error)
	GetAllStores(ctx context.Context, search string) ([]*model.Store, error)
}

type storeService struct {
	repo          Repository
	companyClient company.CompanyServiceClient
}

func NewService(repo Repository, companyClient company.CompanyServiceClient) Service {
	return &storeService{
		repo:          repo,
		companyClient: companyClient,
	}
}

// FIX: Tambahkan kembali fungsi helper yang hilang untuk generate kode unik.
func (s *storeService) generateUniqueStoreCode(ctx context.Context) (string, error) {
	for i := 0; i < 10; i++ {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		code := fmt.Sprintf("%04d", r.Intn(9000)+1000)
		_, err := s.repo.FindByCode(ctx, code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return code, nil
			}
			return "", err
		}
	}
	return "", errors.New("failed to generate a unique store code after 10 attempts")
}

func (s *storeService) CreateStore(ctx context.Context, companyID, name, address string) (*model.Store, error) {
	companyUUID, err := uuid.Parse(companyID)
	if err != nil {
		return nil, errors.New("invalid company id format")
	}

	_, err = s.companyClient.GetCompany(ctx, &company.GetCompanyRequest{Id: companyID})
	if err != nil {
		return nil, errors.New("company with the given ID does not exist")
	}

	code, err := s.generateUniqueStoreCode(ctx)
	if err != nil {
		return nil, err
	}

	store := &model.Store{
		CompanyID: companyUUID,
		Name:      name,
		Address:   address,
		Code:      code,
	}

	if err := s.repo.Create(ctx, store); err != nil {
		return nil, err
	}
	return store, nil
}

// FIX: Implementasikan metode GetStoreByID yang dibutuhkan oleh interface.
func (s *storeService) GetStoreByID(ctx context.Context, id string) (*model.Store, error) {
	storeUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid store id format")
	}
	return s.repo.FindByID(ctx, storeUUID)
}

// FIX: Implementasikan metode GetAllStores yang dibutuhkan oleh interface.
func (s *storeService) GetAllStores(ctx context.Context, search string) ([]*model.Store, error) {
	return s.repo.FindAll(ctx, search)
}
