package model

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}

type UpdateUserRequest struct {
	ID       string  `json:"-"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	FullName *string `json:"full_name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type ListUsersRequest struct {
	Page     int32   `json:"page" form:"page"`
	PageSize int32   `json:"page_size" form:"page_size"`
	SortBy   *string `json:"sort_by,omitempty" form:"sort_by"`
	Filter   *string `json:"filter,omitempty" form:"filter"`
}

type UserResponse struct {
	User       *User  `json:"user,omitempty"`
	Users      []User `json:"users,omitempty"`
	TotalCount int32  `json:"total_count,omitempty"`
	Page       int32  `json:"page,omitempty"`
	PageSize   int32  `json:"page_size,omitempty"`
	Message    string `json:"message,omitempty"`
	Success    bool   `json:"success,omitempty"`
}