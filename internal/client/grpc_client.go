package client

import (
	"context"
	"fmt"

	productpb "go-grpc-rest-demo/api/gen/go/product/v1"
	userpb "go-grpc-rest-demo/api/gen/go/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClient wraps gRPC service clients
type GRPCClient struct {
	conn          *grpc.ClientConn
	userClient    userpb.UserServiceClient
	productClient productpb.ProductServiceClient
	config        *Config
}

// NewGRPCClient creates a new gRPC client
func NewGRPCClient(config *Config) (*GRPCClient, error) {
	conn, err := grpc.NewClient(config.GRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server at %s: %v", config.GRPCAddr, err)
	}

	return &GRPCClient{
		conn:          conn,
		userClient:    userpb.NewUserServiceClient(conn),
		productClient: productpb.NewProductServiceClient(conn),
		config:        config,
	}, nil
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

// User service methods

func (c *GRPCClient) CreateUser(ctx context.Context, username, email, fullName string) (*userpb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &userpb.CreateUserRequest{
		Username: username,
		Email:    email,
		FullName: fullName,
	}

	resp, err := c.userClient.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (c *GRPCClient) GetUser(ctx context.Context, id string) (*userpb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &userpb.GetUserRequest{Id: id}

	resp, err := c.userClient.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (c *GRPCClient) UpdateUser(ctx context.Context, id string, username, email, fullName *string, isActive *bool) (*userpb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &userpb.UpdateUserRequest{
		Id:       id,
		Username: username,
		Email:    email,
		FullName: fullName,
		IsActive: isActive,
	}

	resp, err := c.userClient.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (c *GRPCClient) DeleteUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &userpb.DeleteUserRequest{Id: id}

	_, err := c.userClient.DeleteUser(ctx, req)
	return err
}

func (c *GRPCClient) ListUsers(ctx context.Context, page, pageSize int32, sortBy, filter *string) ([]*userpb.User, int32, int32, int32, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &userpb.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		Filter:   filter,
	}

	resp, err := c.userClient.ListUsers(ctx, req)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return resp.Users, resp.TotalCount, resp.Page, resp.PageSize, nil
}

// Product service methods

func (c *GRPCClient) CreateProduct(ctx context.Context, name, description, category string, price float64, quantity int32) (*productpb.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &productpb.CreateProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		Category:    category,
	}

	resp, err := c.productClient.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Product, nil
}

func (c *GRPCClient) GetProduct(ctx context.Context, id string) (*productpb.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &productpb.GetProductRequest{Id: id}

	resp, err := c.productClient.GetProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Product, nil
}

func (c *GRPCClient) SearchProducts(ctx context.Context, query, category *string, minPrice, maxPrice *float64, page, pageSize int32) ([]*productpb.Product, int32, int32, int32, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &productpb.SearchProductsRequest{
		Query:    query,
		Category: category,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.productClient.SearchProducts(ctx, req)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return resp.Products, resp.TotalCount, resp.Page, resp.PageSize, nil
}
