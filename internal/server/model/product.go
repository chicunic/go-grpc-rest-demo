package model

import (
	"time"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int32     `json:"quantity"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Quantity    int32   `json:"quantity" binding:"required,min=0"`
	Category    string  `json:"category" binding:"required"`
}

type SearchProductsRequest struct {
	Query    *string  `json:"query,omitempty" form:"query"`
	Category *string  `json:"category,omitempty" form:"category"`
	MinPrice *float64 `json:"min_price,omitempty" form:"min_price"`
	MaxPrice *float64 `json:"max_price,omitempty" form:"max_price"`
	Page     int32    `json:"page" form:"page"`
	PageSize int32    `json:"page_size" form:"page_size"`
}

type ProductResponse struct {
	Product    *Product  `json:"product,omitempty"`
	Products   []Product `json:"products,omitempty"`
	TotalCount int32     `json:"total_count,omitempty"`
	Page       int32     `json:"page,omitempty"`
	PageSize   int32     `json:"page_size,omitempty"`
	Message    string    `json:"message,omitempty"`
}