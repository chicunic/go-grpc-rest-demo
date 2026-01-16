package service

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go-grpc-rest-demo/internal/server/errors"
	"go-grpc-rest-demo/internal/server/model"
)

type ProductService struct {
	products map[string]*model.Product
	nextID   int64
	mu       sync.RWMutex
}

func NewProductService() *ProductService {
	return &ProductService{
		products: make(map[string]*model.Product),
		nextID:   1,
	}
}

func (s *ProductService) generateID() string {
	return strconv.FormatInt(atomic.AddInt64(&s.nextID, 1)-1, 10)
}

func (s *ProductService) CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.Product, error) {
	if req.Name == "" || req.Description == "" || req.Category == "" {
		return nil, errors.NewValidationError("fields", "name, description, and category are required")
	}
	if req.Price < 0 || req.Quantity < 0 {
		return nil, errors.NewValidationError("value", "price and quantity cannot be negative")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	product := &model.Product{
		ID:          s.generateID(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Category:    req.Category,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.products[product.ID] = product
	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*model.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, exists := s.products[id]
	if !exists {
		return nil, errors.NewNotFoundError("product", id)
	}
	return product, nil
}

func (s *ProductService) SearchProducts(ctx context.Context, req *model.SearchProductsRequest) ([]model.Product, int32, int32, int32, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := s.filterProducts(req)

	sort.Slice(products, func(i, j int) bool {
		return products[i].Name < products[j].Name
	})

	paged, total, page, pageSize := paginate(products, req.Page, req.PageSize)
	return paged, total, page, pageSize, nil
}

func (s *ProductService) filterProducts(req *model.SearchProductsRequest) []model.Product {
	var products []model.Product
	var queryLower string
	if req.Query != nil {
		queryLower = strings.ToLower(*req.Query)
	}

	for _, product := range s.products {
		if !s.matchesSearchCriteria(product, queryLower, req) {
			continue
		}
		products = append(products, *product)
	}
	return products
}

func (s *ProductService) matchesSearchCriteria(product *model.Product, query string, req *model.SearchProductsRequest) bool {
	if query != "" {
		nameMatch := strings.Contains(strings.ToLower(product.Name), query)
		descMatch := strings.Contains(strings.ToLower(product.Description), query)
		if !nameMatch && !descMatch {
			return false
		}
	}
	if req.Category != nil && !strings.EqualFold(product.Category, *req.Category) {
		return false
	}
	if req.MinPrice != nil && product.Price < *req.MinPrice {
		return false
	}
	if req.MaxPrice != nil && product.Price > *req.MaxPrice {
		return false
	}
	return true
}
