// pkg/grpc/client/client.go

package client

import (
	"fmt"
	"os" // <-- TAMBAHKAN IMPORT OS

	"github.com/ryannovarypradana/fnb-microservice-api/config"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/promotion"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Helper function untuk koneksi gRPC
func dialService(envVarName string) (*grpc.ClientConn, error) {
	addr := os.Getenv(envVarName)
	if addr == "" {
		return nil, fmt.Errorf("environment variable %s not set", envVarName)
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", addr, err)
	}
	return conn, nil
}

func NewAuthClient(cfg *config.Config) (auth.AuthServiceClient, error) {
	conn, err := dialService("AUTH_SERVICE_ADDR") // <-- GUNAKAN ENV VAR
	if err != nil {
		return nil, err
	}
	return auth.NewAuthServiceClient(conn), nil
}

func NewUserClient(cfg *config.Config) (user.UserServiceClient, error) {
	conn, err := dialService("USER_SERVICE_ADDR") // <-- GUNAKAN ENV VAR
	if err != nil {
		return nil, err
	}
	return user.NewUserServiceClient(conn), nil
}

func NewCompanyClient(cfg *config.Config) (company.CompanyServiceClient, error) {
	conn, err := dialService("COMPANY_SERVICE_ADDR") // <-- GUNAKAN ENV VAR
	if err != nil {
		return nil, err
	}
	return company.NewCompanyServiceClient(conn), nil
}

// ... Lakukan hal yang sama untuk semua fungsi New...Client lainnya ...

func NewStoreClient(cfg *config.Config) (store.StoreServiceClient, error) {
	conn, err := dialService("STORE_SERVICE_ADDR")
	if err != nil {
		return nil, err
	}
	return store.NewStoreServiceClient(conn), nil
}

func NewProductClient(cfg *config.Config) (product.ProductServiceClient, error) {
	conn, err := dialService("PRODUCT_SERVICE_ADDR")
	if err != nil {
		return nil, err
	}
	return product.NewProductServiceClient(conn), nil
}

func NewOrderClient(cfg *config.Config) (order.OrderServiceClient, error) {
	conn, err := dialService("ORDER_SERVICE_ADDR")
	if err != nil {
		return nil, err
	}
	return order.NewOrderServiceClient(conn), nil
}

func NewPromotionClient(cfg *config.Config) (promotion.PromotionServiceClient, error) {
	conn, err := dialService("PROMOTION_SERVICE_ADDR")
	if err != nil {
		return nil, err
	}
	return promotion.NewPromotionServiceClient(conn), nil
}
