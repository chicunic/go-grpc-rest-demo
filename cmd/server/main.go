// Package main Go gRPC REST Demo API
//
// This is a demo API server with both gRPC and REST endpoints based on protobuf definitions.
//
// @title Go gRPC REST Demo API
// @version 1.0
// @description This is a demo API server with both gRPC and REST endpoints based on protobuf definitions.
//
// @host localhost:8080
// @BasePath /api/v1
//
// @schemes http
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	productpb "go-grpc-rest-demo/api/gen/go/product/v1"
	userpb "go-grpc-rest-demo/api/gen/go/user/v1"
	_ "go-grpc-rest-demo/docs" // Import docs for swagger
	grpcserver "go-grpc-rest-demo/internal/server/grpc"
	"go-grpc-rest-demo/internal/server/rest"
	"go-grpc-rest-demo/internal/server/service"
)

const (
	restPort = ":8080"
	grpcPort = ":9090"
)

func main() {
	userService := service.NewUserService()
	productService := service.NewProductService()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := runREST(ctx, userService, productService); err != nil {
			log.Printf("REST server error: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := runGRPC(ctx, userService, productService); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	log.Println("Servers started. REST: :8080, gRPC: :9090")
	<-ctx.Done()
	log.Println("Shutting down...")
	wg.Wait()
	log.Println("Shutdown complete")
}

func runREST(ctx context.Context, userService *service.UserService, productService *service.ProductService) error {
	srv := &http.Server{
		Addr:    restPort,
		Handler: rest.SetupRouter(userService, productService),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("REST listen error: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func runGRPC(ctx context.Context, userService *service.UserService, productService *service.ProductService) error {
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, grpcserver.NewUserServer(userService))
	productpb.RegisterProductServiceServer(grpcServer, grpcserver.NewProductServer(productService))
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return err
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC serve error: %v", err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
