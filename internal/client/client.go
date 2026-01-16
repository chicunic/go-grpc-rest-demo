package client

import (
	"context"
	"fmt"

	productpb "go-grpc-rest-demo/api/gen/go/product/v1"
	userpb "go-grpc-rest-demo/api/gen/go/user/v1"
	"go-grpc-rest-demo/internal/server/model"
)

// Client interface defines the methods for both gRPC and REST clients
type Client interface {
	Close() error

	// User methods - return different types based on client type
	CreateUserGRPC(ctx context.Context, username, email, fullName string) (*userpb.User, error)
	CreateUserREST(ctx context.Context, username, email, fullName string) (*model.User, error)

	GetUserGRPC(ctx context.Context, id string) (*userpb.User, error)
	GetUserREST(ctx context.Context, id string) (*model.User, error)

	UpdateUserGRPC(ctx context.Context, id string, username, email, fullName *string, isActive *bool) (*userpb.User, error)
	UpdateUserREST(ctx context.Context, id string, username, email, fullName *string, isActive *bool) (*model.User, error)

	DeleteUser(ctx context.Context, id string) error

	ListUsersGRPC(ctx context.Context, page, pageSize int32, sortBy, filter *string) ([]*userpb.User, int32, int32, int32, error)
	ListUsersREST(ctx context.Context, page, pageSize int32, sortBy, filter *string) ([]model.User, int32, int32, int32, error)

	// Product methods
	CreateProductGRPC(ctx context.Context, name, description, category string, price float64, quantity int32) (*productpb.Product, error)
	CreateProductREST(ctx context.Context, name, description, category string, price float64, quantity int32) (*model.Product, error)

	GetProductGRPC(ctx context.Context, id string) (*productpb.Product, error)
	GetProductREST(ctx context.Context, id string) (*model.Product, error)

	SearchProductsGRPC(ctx context.Context, query, category *string, minPrice, maxPrice *float64, page, pageSize int32) ([]*productpb.Product, int32, int32, int32, error)
	SearchProductsREST(ctx context.Context, query, category *string, minPrice, maxPrice *float64, page, pageSize int32) ([]model.Product, int32, int32, int32, error)
}

// UnifiedClient wraps both gRPC and REST clients
type UnifiedClient struct {
	grpcClient *GRPCClient
	restClient *RESTClient
	config     *Config
}

// NewClient creates a new unified client based on configuration
func NewClient(config *Config) (Client, error) {
	var grpcClient *GRPCClient
	var restClient *RESTClient
	var err error

	// Always try to create both clients, but only require the one specified in mode
	if config.Mode == "grpc" || config.Mode == "both" {
		grpcClient, err = NewGRPCClient(config)
		if err != nil && config.Mode == "grpc" {
			return nil, fmt.Errorf("failed to create gRPC client: %v", err)
		}
	}

	if config.Mode == "rest" || config.Mode == "both" {
		restClient, err = NewRESTClient(config)
		if err != nil && config.Mode == "rest" {
			return nil, fmt.Errorf("failed to create REST client: %v", err)
		}
	}

	return &UnifiedClient{
		grpcClient: grpcClient,
		restClient: restClient,
		config:     config,
	}, nil
}

func (c *UnifiedClient) Close() error {
	var err error
	if c.grpcClient != nil {
		if grpcErr := c.grpcClient.Close(); grpcErr != nil {
			err = grpcErr
		}
	}
	if c.restClient != nil {
		if restErr := c.restClient.Close(); restErr != nil {
			err = restErr
		}
	}
	return err
}

// gRPC methods
func (c *UnifiedClient) CreateUserGRPC(ctx context.Context, username, email, fullName string) (*userpb.User, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client not available")
	}
	return c.grpcClient.CreateUser(ctx, username, email, fullName)
}

func (c *UnifiedClient) GetUserGRPC(ctx context.Context, id string) (*userpb.User, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client not available")
	}
	return c.grpcClient.GetUser(ctx, id)
}

func (c *UnifiedClient) UpdateUserGRPC(ctx context.Context, id string, username, email, fullName *string, isActive *bool) (*userpb.User, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client not available")
	}
	return c.grpcClient.UpdateUser(ctx, id, username, email, fullName, isActive)
}

func (c *UnifiedClient) ListUsersGRPC(ctx context.Context, page, pageSize int32, sortBy, filter *string) ([]*userpb.User, int32, int32, int32, error) {
	if c.grpcClient == nil {
		return nil, 0, 0, 0, fmt.Errorf("gRPC client not available")
	}
	return c.grpcClient.ListUsers(ctx, page, pageSize, sortBy, filter)
}

func (c *UnifiedClient) CreateProductGRPC(ctx context.Context, name, description, category string, price float64, quantity int32) (*productpb.Product, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client not available")
	}
	return c.grpcClient.CreateProduct(ctx, name, description, category, price, quantity)
}

func (c *UnifiedClient) GetProductGRPC(ctx context.Context, id string) (*productpb.Product, error) {
	if c.grpcClient == nil {
		return nil, fmt.Errorf("gRPC client not available")
	}
	return c.grpcClient.GetProduct(ctx, id)
}

func (c *UnifiedClient) SearchProductsGRPC(ctx context.Context, query, category *string, minPrice, maxPrice *float64, page, pageSize int32) ([]*productpb.Product, int32, int32, int32, error) {
	if c.grpcClient == nil {
		return nil, 0, 0, 0, fmt.Errorf("gRPC client not available")
	}
	return c.grpcClient.SearchProducts(ctx, query, category, minPrice, maxPrice, page, pageSize)
}

// REST methods
func (c *UnifiedClient) CreateUserREST(ctx context.Context, username, email, fullName string) (*model.User, error) {
	if c.restClient == nil {
		return nil, fmt.Errorf("REST client not available")
	}
	return c.restClient.CreateUser(ctx, username, email, fullName)
}

func (c *UnifiedClient) GetUserREST(ctx context.Context, id string) (*model.User, error) {
	if c.restClient == nil {
		return nil, fmt.Errorf("REST client not available")
	}
	return c.restClient.GetUser(ctx, id)
}

func (c *UnifiedClient) UpdateUserREST(ctx context.Context, id string, username, email, fullName *string, isActive *bool) (*model.User, error) {
	if c.restClient == nil {
		return nil, fmt.Errorf("REST client not available")
	}
	return c.restClient.UpdateUser(ctx, id, username, email, fullName, isActive)
}

func (c *UnifiedClient) ListUsersREST(ctx context.Context, page, pageSize int32, sortBy, filter *string) ([]model.User, int32, int32, int32, error) {
	if c.restClient == nil {
		return nil, 0, 0, 0, fmt.Errorf("REST client not available")
	}
	return c.restClient.ListUsers(ctx, page, pageSize, sortBy, filter)
}

func (c *UnifiedClient) CreateProductREST(ctx context.Context, name, description, category string, price float64, quantity int32) (*model.Product, error) {
	if c.restClient == nil {
		return nil, fmt.Errorf("REST client not available")
	}
	return c.restClient.CreateProduct(ctx, name, description, category, price, quantity)
}

func (c *UnifiedClient) GetProductREST(ctx context.Context, id string) (*model.Product, error) {
	if c.restClient == nil {
		return nil, fmt.Errorf("REST client not available")
	}
	return c.restClient.GetProduct(ctx, id)
}

func (c *UnifiedClient) SearchProductsREST(ctx context.Context, query, category *string, minPrice, maxPrice *float64, page, pageSize int32) ([]model.Product, int32, int32, int32, error) {
	if c.restClient == nil {
		return nil, 0, 0, 0, fmt.Errorf("REST client not available")
	}
	return c.restClient.SearchProducts(ctx, query, category, minPrice, maxPrice, page, pageSize)
}

// Shared methods (work for both)
func (c *UnifiedClient) DeleteUser(ctx context.Context, id string) error {
	if c.config.Mode == "grpc" && c.grpcClient != nil {
		return c.grpcClient.DeleteUser(ctx, id)
	} else if c.config.Mode == "rest" && c.restClient != nil {
		return c.restClient.DeleteUser(ctx, id)
	}
	return fmt.Errorf("no client available for mode: %s", c.config.Mode)
}
