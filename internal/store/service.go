package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	pb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"gorm.io/gorm"
)

type StoreService interface {
	CreateStore(ctx context.Context, req *pb.CreateStoreRequest) (*model.Store, error)
	GetStore(ctx context.Context, id string) (*model.Store, error)
	GetAllStores(ctx context.Context, search string) ([]*model.Store, error)
	UpdateStore(ctx context.Context, req *pb.UpdateStoreRequest) (*model.Store, error)
	DeleteStore(ctx context.Context, id string) error
	CloneStoreContent(ctx context.Context, sourceStoreID, destStoreID string) error
	GetStoreByCode(ctx context.Context, code string) (*model.Store, error)
}

type storeService struct {
	repo StoreRepository
}

func NewStoreService(repo StoreRepository) StoreService {
	return &storeService{repo: repo}
}

// CreateStore adalah handler untuk membuat toko baru.
func (s *storeService) CreateStore(ctx context.Context, req *pb.CreateStoreRequest) (*model.Store, error) {
	companyID, err := uuid.Parse(req.CompanyId)
	if err != nil {
		return nil, errors.New("invalid company id format")
	}

	// Memanggil fungsi yang diperbarui untuk memastikan kode unik
	uniqueCode, err := s.generateUniqueStoreCode(ctx)
	if err != nil {
		return nil, err
	}

	store := &model.Store{
		Name:      req.Name,
		Location:  req.Address,
		CompanyID: companyID,
		Code:      uniqueCode, // Gunakan kode yang unik
	}

	if err := s.repo.Create(ctx, store); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *storeService) GetStore(ctx context.Context, id string) (*model.Store, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *storeService) GetAllStores(ctx context.Context, search string) ([]*model.Store, error) {
	return s.repo.FindAll(ctx, search)
}

func (s *storeService) UpdateStore(ctx context.Context, req *pb.UpdateStoreRequest) (*model.Store, error) {
	store, err := s.repo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		store.Name = *req.Name
	}
	if req.Address != nil {
		store.Location = *req.Address
	}
	if req.BannerImageUrl != nil {
		store.BannerImageURL = *req.BannerImageUrl
	}
	if req.TaxPercentage != nil {
		tax := float64(*req.TaxPercentage)
		store.TaxPercentage = &tax
	}
	if req.Latitude != nil {
		lat := float64(*req.Latitude)
		store.Latitude = &lat
	}
	if req.Longitude != nil {
		long := float64(*req.Longitude)
		store.Longitude = &long
	}
	if req.OperationalHours != nil {
		var opHours model.OpeningHours
		if err := json.Unmarshal([]byte(*req.OperationalHours), &opHours); err == nil {
			store.OperationalHours = &opHours
		} else {
			log.Printf("Warning: could not unmarshal operational hours: %v", err)
		}
	}

	if err := s.repo.Update(ctx, store); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *storeService) DeleteStore(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *storeService) CloneStoreContent(ctx context.Context, sourceStoreID, destStoreID string) error {
	log.Printf("Placeholder: Cloning content from store %s to %s", sourceStoreID, destStoreID)
	return nil
}

func (s *storeService) GetStoreByCode(ctx context.Context, code string) (*model.Store, error) {
	return s.repo.FindByCode(ctx, code)
}

// generateUniqueStoreCode adalah helper function yang diperbarui
func (s *storeService) generateUniqueStoreCode(ctx context.Context) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		code := fmt.Sprintf("%04d", r.Intn(9000)+1000)
		_, err := s.repo.FindByCode(ctx, code)
		if errors.Is(err, gorm.ErrRecordNotFound) { // Periksa error GORM yang asli
			return code, nil
		}
		if err != nil {
			// Jika ada error lain (misal: koneksi DB), kembalikan error tersebut
			return "", err
		}
	}
	return "", errors.New("gagal membuat kode toko unik setelah 100 percobaan")
}
