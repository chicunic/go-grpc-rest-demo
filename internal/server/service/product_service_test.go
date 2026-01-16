package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-grpc-rest-demo/internal/server/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProductServiceTestSuite struct {
	suite.Suite
	service *ProductService
}

func (suite *ProductServiceTestSuite) SetupTest() {
	suite.service = NewProductService()
}

func (suite *ProductServiceTestSuite) TestCreateProduct() {
	req := &model.CreateProductRequest{
		Name:        "Test Product",
		Description: "A test product",
		Price:       19.99,
		Quantity:    100,
		Category:    "Electronics",
	}

	product, err := suite.service.CreateProduct(context.Background(), req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), product)
	assert.Equal(suite.T(), req.Name, product.Name)
	assert.Equal(suite.T(), req.Description, product.Description)
	assert.Equal(suite.T(), req.Price, product.Price)
	assert.Equal(suite.T(), req.Quantity, product.Quantity)
	assert.Equal(suite.T(), req.Category, product.Category)
	assert.NotEmpty(suite.T(), product.ID)
	assert.WithinDuration(suite.T(), time.Now(), product.CreatedAt, time.Second)
	assert.WithinDuration(suite.T(), time.Now(), product.UpdatedAt, time.Second)
}

func (suite *ProductServiceTestSuite) TestGetProduct() {
	req := &model.CreateProductRequest{
		Name:        "Get Product",
		Description: "A product to get",
		Price:       29.99,
		Quantity:    50,
		Category:    "Books",
	}

	created, err := suite.service.CreateProduct(context.Background(), req)
	assert.NoError(suite.T(), err)

	retrieved, err := suite.service.GetProduct(context.Background(), created.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), created.ID, retrieved.ID)
	assert.Equal(suite.T(), created.Name, retrieved.Name)
	assert.Equal(suite.T(), created.Description, retrieved.Description)
	assert.Equal(suite.T(), created.Price, retrieved.Price)
	assert.Equal(suite.T(), created.Quantity, retrieved.Quantity)
	assert.Equal(suite.T(), created.Category, retrieved.Category)
}

func (suite *ProductServiceTestSuite) TestGetProductNotFound() {
	_, err := suite.service.GetProduct(context.Background(), "nonexistent")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *ProductServiceTestSuite) TestSearchProducts() {
	categories := []string{"Electronics", "Books", "Electronics", "Clothing", "Books"}

	for i := 0; i < 5; i++ {
		createReq := &model.CreateProductRequest{
			Name:        fmt.Sprintf("Product %d", i),
			Description: fmt.Sprintf("Description %d", i),
			Price:       float64(10 + i*5),
			Quantity:    int32(50 + i*10),
			Category:    categories[i],
		}
		_, err := suite.service.CreateProduct(context.Background(), createReq)
		assert.NoError(suite.T(), err)
	}

	// Test basic pagination
	searchReq := &model.SearchProductsRequest{
		Page:     1,
		PageSize: 3,
	}
	result, totalCount, page, pageSize, err := suite.service.SearchProducts(context.Background(), searchReq)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)
	assert.Equal(suite.T(), int32(5), totalCount)
	assert.Equal(suite.T(), int32(1), page)
	assert.Equal(suite.T(), int32(3), pageSize)

	// Test category filter
	categoryFilter := "Electronics"
	searchReq = &model.SearchProductsRequest{
		Category: &categoryFilter,
		Page:     1,
		PageSize: 10,
	}
	result, totalCount, _, _, err = suite.service.SearchProducts(context.Background(), searchReq)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), int32(2), totalCount)
	for _, product := range result {
		assert.Equal(suite.T(), "Electronics", product.Category)
	}

	// Test query filter
	query := "Product 1"
	searchReq = &model.SearchProductsRequest{
		Query:    &query,
		Page:     1,
		PageSize: 10,
	}
	result, totalCount, _, _, err = suite.service.SearchProducts(context.Background(), searchReq)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), int32(1), totalCount)
	assert.Contains(suite.T(), result[0].Name, "1")

	// Test price range filter
	minPrice := 15.0
	maxPrice := 25.0
	searchReq = &model.SearchProductsRequest{
		MinPrice: &minPrice,
		MaxPrice: &maxPrice,
		Page:     1,
		PageSize: 10,
	}
	result, totalCount, _, _, err = suite.service.SearchProducts(context.Background(), searchReq)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), totalCount >= 1)
	for _, product := range result {
		assert.True(suite.T(), product.Price >= minPrice)
		assert.True(suite.T(), product.Price <= maxPrice)
	}
}

func (suite *ProductServiceTestSuite) TestSearchProductsEmpty() {
	// Test search with no products
	searchReq := &model.SearchProductsRequest{
		Page:     1,
		PageSize: 10,
	}
	result, totalCount, page, pageSize, err := suite.service.SearchProducts(context.Background(), searchReq)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
	assert.Equal(suite.T(), int32(0), totalCount)
	assert.Equal(suite.T(), int32(1), page)
	assert.Equal(suite.T(), int32(10), pageSize)
}

func (suite *ProductServiceTestSuite) TestSearchProductsWithNonexistentCategory() {
	createReq := &model.CreateProductRequest{
		Name:        "Test Product",
		Description: "Description",
		Price:       19.99,
		Quantity:    100,
		Category:    "Electronics",
	}
	_, err := suite.service.CreateProduct(context.Background(), createReq)
	assert.NoError(suite.T(), err)

	categoryFilter := "Nonexistent"
	searchReq := &model.SearchProductsRequest{
		Category: &categoryFilter,
		Page:     1,
		PageSize: 10,
	}
	result, totalCount, _, _, err := suite.service.SearchProducts(context.Background(), searchReq)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 0)
	assert.Equal(suite.T(), int32(0), totalCount)
}

func TestProductServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}