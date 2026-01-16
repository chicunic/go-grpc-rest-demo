package rest

import (
	"net/http"
	"strconv"

	"go-grpc-rest-demo/internal/server/errors"
	"go-grpc-rest-demo/internal/server/model"
	"go-grpc-rest-demo/internal/server/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the provided information
// @Tags products
// @Accept json
// @Produce json
// @Param product body model.CreateProductRequest true "Product information"
// @Success 201 {object} model.ProductResponse
// @Failure 400 {object} model.ProductResponse
// @Failure 500 {object} model.ProductResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleProductError(c, errors.NewInvalidRequestError("Invalid request: "+err.Error()))
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		handleProductError(c, err)
		return
	}

	respondProductSuccess(c, http.StatusCreated, product)
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} model.ProductResponse
// @Failure 400 {object} model.ProductResponse
// @Failure 404 {object} model.ProductResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleProductError(c, errors.NewValidationError("id", "Product ID is required"))
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), id)
	if err != nil {
		handleProductError(c, err)
		return
	}

	respondProductSuccess(c, http.StatusOK, product)
}

// SearchProducts godoc
// @Summary Search products
// @Description Search products with optional filters
// @Tags products
// @Produce json
// @Param query query string false "Search query (matches name or description)"
// @Param category query string false "Filter by category"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Success 200 {object} model.ProductResponse
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

	req := &model.SearchProductsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	if query := c.Query("query"); query != "" {
		req.Query = &query
	}

	if category := c.Query("category"); category != "" {
		req.Category = &category
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			req.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			req.MaxPrice = &maxPrice
		}
	}

	products, totalCount, retPage, retPageSize, err := h.productService.SearchProducts(c.Request.Context(), req)
	if err != nil {
		handleProductError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.ProductResponse{
		Products:   products,
		TotalCount: totalCount,
		Page:       retPage,
		PageSize:   retPageSize,
		Message:    "Products retrieved successfully",
	})
}
