package client

import (
	"log"

	authPb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	companyPb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
	orderPb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	productPb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	storePb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	userPb "github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// --- Auth Service Client ---
type IAuthServiceClient interface {
	GetAuthServiceClient() authPb.AuthServiceClient
}
type AuthServiceClient struct {
	Client authPb.AuthServiceClient
}

func NewAuthServiceClient(port string) IAuthServiceClient {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to auth service: %v", err)
	}
	return &AuthServiceClient{Client: authPb.NewAuthServiceClient(conn)}
}
func (c *AuthServiceClient) GetAuthServiceClient() authPb.AuthServiceClient {
	return c.Client
}

// --- Product Service Client ---
type IProductServiceClient interface {
	GetProductServiceClient() productPb.ProductServiceClient
}
type ProductServiceClient struct {
	Client productPb.ProductServiceClient
}

func NewProductServiceClient(port string) IProductServiceClient {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to product service: %v", err)
	}
	return &ProductServiceClient{Client: productPb.NewProductServiceClient(conn)}
}
func (p *ProductServiceClient) GetProductServiceClient() productPb.ProductServiceClient {
	return p.Client
}

// --- Order Service Client (Perbaikan Kunci di Sini) ---
type IOrderServiceClient interface {
	GetOrderServiceClient() orderPb.OrderServiceClient
}
type OrderServiceClient struct {
	Client orderPb.OrderServiceClient
}

func NewOrderServiceClient(port string) IOrderServiceClient {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to order service: %v", err)
	}
	return &OrderServiceClient{Client: orderPb.NewOrderServiceClient(conn)}
}
func (o *OrderServiceClient) GetOrderServiceClient() orderPb.OrderServiceClient {
	return o.Client
}

// --- User Service Client ---
type IUserServiceClient interface {
	GetUserServiceClient() userPb.UserServiceClient
}
type UserServiceClient struct {
	Client userPb.UserServiceClient
}

func NewUserServiceClient(port string) IUserServiceClient {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to user service: %v", err)
	}
	return &UserServiceClient{Client: userPb.NewUserServiceClient(conn)}
}
func (u *UserServiceClient) GetUserServiceClient() userPb.UserServiceClient {
	return u.Client
}

// --- Store Service Client ---
type IStoreServiceClient interface {
	GetStoreServiceClient() storePb.StoreServiceClient
}
type StoreServiceClient struct {
	Client storePb.StoreServiceClient
}

func NewStoreServiceClient(port string) IStoreServiceClient {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to store service: %v", err)
	}
	return &StoreServiceClient{Client: storePb.NewStoreServiceClient(conn)}
}
func (s *StoreServiceClient) GetStoreServiceClient() storePb.StoreServiceClient {
	return s.Client
}

// --- Company Service Client ---
type ICompanyServiceClient interface {
	GetCompanyServiceClient() companyPb.CompanyServiceClient
}
type CompanyServiceClient struct {
	Client companyPb.CompanyServiceClient
}

func NewCompanyServiceClient(port string) ICompanyServiceClient {
	conn, err := grpc.Dial("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to company service: %v", err)
	}
	return &CompanyServiceClient{Client: companyPb.NewCompanyServiceClient(conn)}
}
func (c *CompanyServiceClient) GetCompanyServiceClient() companyPb.CompanyServiceClient {
	return c.Client
}
