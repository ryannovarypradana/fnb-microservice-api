// pkg/grpc/client/client.go

package client

import (
	"fmt"

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"   // <-- IMPORT DITAMBAHKAN
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product" // <-- IMPORT DITAMBAHKAN
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(cfg *config.Config) (auth.AuthServiceClient, error) {
	addr := fmt.Sprintf("localhost:%s", cfg.Auth.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial auth service: %w", err)
	}
	return auth.NewAuthServiceClient(conn), nil
}

func NewUserClient(cfg *config.Config) (user.UserServiceClient, error) {
	addr := fmt.Sprintf("localhost:%s", cfg.User.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial user service: %w", err)
	}
	return user.NewUserServiceClient(conn), nil
}

func NewCompanyClient(cfg *config.Config) (company.CompanyServiceClient, error) {
	addr := fmt.Sprintf("localhost:%s", cfg.Company.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial company service: %w", err)
	}
	return company.NewCompanyServiceClient(conn), nil
}

func NewStoreClient(cfg *config.Config) (store.StoreServiceClient, error) {
	addr := fmt.Sprintf("localhost:%s", cfg.Store.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial store service: %w", err)
	}
	return store.NewStoreServiceClient(conn), nil
}

func NewProductClient(cfg *config.Config) (product.ProductServiceClient, error) {
	addr := fmt.Sprintf("localhost:%s", cfg.Product.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial product service: %w", err)
	}
	return product.NewProductServiceClient(conn), nil
}

func NewOrderClient(cfg *config.Config) (order.OrderServiceClient, error) {
	addr := fmt.Sprintf("localhost:%s", cfg.Order.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial order service: %w", err)
	}
	return order.NewOrderServiceClient(conn), nil
}
