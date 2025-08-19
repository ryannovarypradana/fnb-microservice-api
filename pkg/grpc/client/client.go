package client

import (
	"log"
	"os"

	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/auth"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/company"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/order"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/product"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/store"
	"github.com/ryannovarypradana/fnb-microservice-api/pkg/grpc/protoc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newGRPCConnection(target string) *grpc.ClientConn {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server at %s: %v", target, err)
	}
	return conn
}

func NewAuthClient() auth.AuthServiceClient {
	target := os.Getenv("AUTH_SERVICE_URL") // e.g., "localhost:50051"
	conn := newGRPCConnection(target)
	return auth.NewAuthServiceClient(conn)
}

func NewUserClient() user.UserServiceClient {
	target := os.Getenv("USER_SERVICE_URL") // e.g., "localhost:50052"
	conn := newGRPCConnection(target)
	return user.NewUserServiceClient(conn)
}

func NewProductClient() product.ProductServiceClient {
	target := os.Getenv("PRODUCT_SERVICE_URL") // e.g., "localhost:50053"
	conn := newGRPCConnection(target)
	return product.NewProductServiceClient(conn)
}

func NewCompanyClient() company.CompanyServiceClient {
	target := os.Getenv("COMPANY_SERVICE_URL") // e.g., "localhost:50054"
	conn := newGRPCConnection(target)
	return company.NewCompanyServiceClient(conn)
}

func NewStoreClient() store.StoreServiceClient {
	target := os.Getenv("STORE_SERVICE_URL") // e.g., "localhost:50055"
	conn := newGRPCConnection(target)
	return store.NewStoreServiceClient(conn)
}

func NewOrderClient() order.OrderServiceClient {
	target := os.Getenv("ORDER_SERVICE_URL") // e.g., "localhost:50056"
	conn := newGRPCConnection(target)
	return order.NewOrderServiceClient(conn)
}
