package grpc

import (
	"context"
	"time"

	pb "go-grpc-rest-demo/api/gen/go/product/v1"
	"go-grpc-rest-demo/internal/server/errors"
	"go-grpc-rest-demo/internal/server/model"
	"go-grpc-rest-demo/internal/server/service"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
	productService *service.ProductService
}

func NewProductServer(productService *service.ProductService) *ProductServer {
	return &ProductServer{productService: productService}
}

func productToPB(product *model.Product) *pb.Product {
	return &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	if req.Name == "" || req.Description == "" || req.Category == "" {
		return nil, handleGRPCError(errors.NewValidationError("fields", "name, description, and category are required"))
	}
	if req.Price < 0 {
		return nil, handleGRPCError(errors.NewValidationError("price", "price must be non-negative"))
	}
	if req.Quantity < 0 {
		return nil, handleGRPCError(errors.NewValidationError("quantity", "quantity must be non-negative"))
	}

	modelReq := &model.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		Category:    req.Category,
	}

	product, err := s.productService.CreateProduct(ctx, modelReq)
	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.CreateProductResponse{
		Product: productToPB(product),
		Message: "Product created successfully",
	}, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	if req.Id == "" {
		return nil, handleGRPCError(errors.NewValidationError("id", "id is required"))
	}

	product, err := s.productService.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.GetProductResponse{
		Product: productToPB(product),
		Message: "Product retrieved successfully",
	}, nil
}

func (s *ProductServer) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	modelReq := &model.SearchProductsRequest{
		Query:    req.Query,
		Category: req.Category,
		MinPrice: req.MinPrice,
		MaxPrice: req.MaxPrice,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	products, totalCount, page, pageSize, err := s.productService.SearchProducts(ctx, modelReq)
	if err != nil {
		return nil, handleGRPCError(err)
	}

	pbProducts := make([]*pb.Product, len(products))
	for i := range products {
		pbProducts[i] = productToPB(&products[i])
	}

	return &pb.SearchProductsResponse{
		Products:   pbProducts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
