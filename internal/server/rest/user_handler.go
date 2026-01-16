package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-grpc-rest-demo/internal/server/errors"
	"go-grpc-rest-demo/internal/server/model"
	"go-grpc-rest-demo/internal/server/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.CreateUserRequest true "User information"
// @Success 201 {object} model.UserResponse
// @Failure 400 {object} model.UserResponse
// @Failure 500 {object} model.UserResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := errors.NewInvalidRequestError("Invalid JSON format: " + err.Error())
		handleUserError(c, appErr)
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		handleUserError(c, err)
		return
	}

	respondUserSuccess(c, http.StatusCreated, user)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} model.UserResponse
// @Failure 404 {object} model.UserResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleUserError(c, errors.NewValidationError("id", "User ID is required"))
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		handleUserError(c, err)
		return
	}

	respondUserSuccess(c, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body model.UpdateUserRequest true "Updated user information"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} model.UserResponse
// @Failure 404 {object} model.UserResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleUserError(c, errors.NewValidationError("id", "User ID is required"))
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleUserError(c, errors.NewInvalidRequestError("Invalid request: "+err.Error()))
		return
	}

	req.ID = id
	user, err := h.userService.UpdateUser(c.Request.Context(), &req)
	if err != nil {
		handleUserError(c, err)
		return
	}

	respondUserSuccess(c, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.UserResponse
// @Failure 400 {object} model.UserResponse
// @Failure 404 {object} model.UserResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleUserError(c, errors.NewValidationError("id", "User ID is required"))
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		handleUserError(c, err)
		return
	}

	respondUserSuccess(c, http.StatusOK, nil)
}

// ListUsers godoc
// @Summary List users
// @Description Get a paginated list of users with optional filtering and sorting
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Items per page" default(10)
// @Param sort_by query string false "Sort by field (username, email, full_name, created_at)"
// @Param filter query string false "Filter by username, email, or full_name"
// @Success 200 {object} model.UserResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 32)

	req := &model.ListUsersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		req.SortBy = &sortBy
	}

	if filter := c.Query("filter"); filter != "" {
		req.Filter = &filter
	}

	users, totalCount, retPage, retPageSize, err := h.userService.ListUsers(c.Request.Context(), req)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, model.UserResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       retPage,
		PageSize:   retPageSize,
		Message:    "Users retrieved successfully",
	})
}