package grpc

import (
	"context"
	"time"

	pb "go-grpc-rest-demo/api/gen/go/user/v1"
	"go-grpc-rest-demo/internal/server/errors"
	"go-grpc-rest-demo/internal/server/model"
	"go-grpc-rest-demo/internal/server/service"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserServer(userService *service.UserService) *UserServer {
	return &UserServer{userService: userService}
}

func userToPB(user *model.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	modelReq := &model.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
	}

	user, err := s.userService.CreateUser(ctx, modelReq)
	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.CreateUserResponse{
		User:    userToPB(user),
		Message: "User created successfully",
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id == "" {
		return nil, handleGRPCError(errors.NewValidationError("id", "id is required"))
	}

	user, err := s.userService.GetUser(ctx, req.Id)
	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.GetUserResponse{
		User:    userToPB(user),
		Message: "User retrieved successfully",
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Id == "" {
		return nil, handleGRPCError(errors.NewValidationError("id", "id is required"))
	}

	modelReq := &model.UpdateUserRequest{
		ID:       req.Id,
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		IsActive: req.IsActive,
	}

	user, err := s.userService.UpdateUser(ctx, modelReq)
	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.UpdateUserResponse{
		User:    userToPB(user),
		Message: "User updated successfully",
	}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id == "" {
		return nil, handleGRPCError(errors.NewValidationError("id", "id is required"))
	}

	if err := s.userService.DeleteUser(ctx, req.Id); err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	modelReq := &model.ListUsersRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		SortBy:   req.SortBy,
		Filter:   req.Filter,
	}

	users, totalCount, page, pageSize, err := s.userService.ListUsers(ctx, modelReq)
	if err != nil {
		return nil, handleGRPCError(err)
	}

	pbUsers := make([]*pb.User, len(users))
	for i := range users {
		pbUsers[i] = userToPB(&users[i])
	}

	return &pb.ListUsersResponse{
		Users:      pbUsers,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
