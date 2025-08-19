package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	// FIX: Import paket eventbus tanpa alias yang salah
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/eventbus"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	RegisterStaff(ctx context.Context, req *user.RegisterStaffRequest) (*model.User, error)
}

type userService struct {
	repo        Repository
	storeClient store.StoreServiceClient
	// FIX: Gunakan tipe yang benar dari paket yang diimpor
	eventPublisher eventbus.Publisher
}

func NewService(repo Repository, storeClient store.StoreServiceClient, publisher eventbus.Publisher) Service {
	return &userService{
		repo:           repo,
		storeClient:    storeClient,
		eventPublisher: publisher,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, userUUID)
}

func (s *userService) RegisterStaff(ctx context.Context, req *user.RegisterStaffRequest) (*model.User, error) {
	storeInfo, err := s.storeClient.GetStore(ctx, &store.GetStoreRequest{Id: req.GetStoreId()})
	if err != nil {
		return nil, errors.New("toko yang dituju tidak valid atau tidak ditemukan")
	}

	if storeInfo.GetStore().GetCompanyId() != req.GetActorCompanyId() {
		return nil, errors.New("akses ditolak: anda tidak bisa mendaftarkan staf untuk toko di luar perusahaan anda")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	storeUUID, _ := uuid.Parse(req.GetStoreId())
	companyUUID, _ := uuid.Parse(req.GetActorCompanyId())

	newUser := &model.User{
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		Password:  string(hashedPassword),
		Role:      req.GetRole(),
		StoreID:   &storeUUID,
		CompanyID: &companyUUID,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("gagal membuat user: %w", err)
	}

	log.Printf("Mempublikasikan event 'user.registered' untuk user %s", newUser.Email)
	eventPayload := map[string]interface{}{
		"user_id": newUser.ID.String(),
		"email":   newUser.Email,
		"name":    newUser.Name,
	}

	go func() {
		err := s.eventPublisher.Publish("user_events", "user.registered", "application/json", eventPayload)
		if err != nil {
			log.Printf("ERROR: Gagal mempublikasikan event user.registered untuk user %s: %v", newUser.Email, err)
		}
	}()

	return newUser, nil
}
