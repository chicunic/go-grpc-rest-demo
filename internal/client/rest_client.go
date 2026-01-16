package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"go-grpc-rest-demo/internal/server/model"
)

// RESTClient wraps HTTP client for REST API calls
type RESTClient struct {
	client  *http.Client
	baseURL string
	config  *Config
}

// NewRESTClient creates a new REST client
func NewRESTClient(config *Config) (*RESTClient, error) {
	return &RESTClient{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		baseURL: config.RESTAddr,
		config:  config,
	}, nil
}

// Close is a no-op for REST client
func (c *RESTClient) Close() error {
	return nil
}

// Helper methods

func (c *RESTClient) doRequest(ctx context.Context, method, path string, body any, result any) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %v", err)
		}
	}

	return nil
}

// User service methods

func (c *RESTClient) CreateUser(ctx context.Context, username, email, fullName string) (*model.User, error) {
	req := &model.CreateUserRequest{
		Username: username,
		Email:    email,
		FullName: fullName,
	}

	var result struct {
		User *model.User `json:"user"`
	}

	err := c.doRequest(ctx, "POST", "/api/v1/users", req, &result)
	if err != nil {
		return nil, err
	}

	return result.User, nil
}

func (c *RESTClient) GetUser(ctx context.Context, id string) (*model.User, error) {
	var result struct {
		User *model.User `json:"user"`
	}

	err := c.doRequest(ctx, "GET", "/api/v1/users/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return result.User, nil
}

func (c *RESTClient) UpdateUser(ctx context.Context, id string, username, email, fullName *string, isActive *bool) (*model.User, error) {
	req := &model.UpdateUserRequest{
		Username: username,
		Email:    email,
		FullName: fullName,
		IsActive: isActive,
	}

	var result struct {
		User *model.User `json:"user"`
	}

	err := c.doRequest(ctx, "PUT", "/api/v1/users/"+id, req, &result)
	if err != nil {
		return nil, err
	}

	return result.User, nil
}

func (c *RESTClient) DeleteUser(ctx context.Context, id string) error {
	return c.doRequest(ctx, "DELETE", "/api/v1/users/"+id, nil, nil)
}

func (c *RESTClient) ListUsers(ctx context.Context, page, pageSize int32, sortBy, filter *string) ([]model.User, int32, int32, int32, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(int(page)))
	params.Set("page_size", strconv.Itoa(int(pageSize)))
	if sortBy != nil {
		params.Set("sort_by", *sortBy)
	}
	if filter != nil {
		params.Set("filter", *filter)
	}

	var result struct {
		Users      []model.User `json:"users"`
		TotalCount int32        `json:"total_count"`
		Page       int32        `json:"page"`
		PageSize   int32        `json:"page_size"`
	}

	err := c.doRequest(ctx, "GET", "/api/v1/users?"+params.Encode(), nil, &result)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return result.Users, result.TotalCount, result.Page, result.PageSize, nil
}

// Product service methods

func (c *RESTClient) CreateProduct(ctx context.Context, name, description, category string, price float64, quantity int32) (*model.Product, error) {
	req := &model.CreateProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		Category:    category,
	}

	var result struct {
		Product *model.Product `json:"product"`
	}

	err := c.doRequest(ctx, "POST", "/api/v1/products", req, &result)
	if err != nil {
		return nil, err
	}

	return result.Product, nil
}

func (c *RESTClient) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	var result struct {
		Product *model.Product `json:"product"`
	}

	err := c.doRequest(ctx, "GET", "/api/v1/products/"+id, nil, &result)
	if err != nil {
		return nil, err
	}

	return result.Product, nil
}

func (c *RESTClient) SearchProducts(ctx context.Context, query, category *string, minPrice, maxPrice *float64, page, pageSize int32) ([]model.Product, int32, int32, int32, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(int(page)))
	params.Set("page_size", strconv.Itoa(int(pageSize)))
	if query != nil {
		params.Set("query", *query)
	}
	if category != nil {
		params.Set("category", *category)
	}
	if minPrice != nil {
		params.Set("min_price", strconv.FormatFloat(*minPrice, 'f', 2, 64))
	}
	if maxPrice != nil {
		params.Set("max_price", strconv.FormatFloat(*maxPrice, 'f', 2, 64))
	}

	var result struct {
		Products   []model.Product `json:"products"`
		TotalCount int32           `json:"total_count"`
		Page       int32           `json:"page"`
		PageSize   int32           `json:"page_size"`
	}

	err := c.doRequest(ctx, "GET", "/api/v1/products/search?"+params.Encode(), nil, &result)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return result.Products, result.TotalCount, result.Page, result.PageSize, nil
}
