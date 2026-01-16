package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-grpc-rest-demo/internal/server/model"
	"go-grpc-rest-demo/internal/server/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	userService *service.UserService
}

func (suite *UserHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.userService = service.NewUserService()
	userHandler := NewUserHandler(suite.userService)

	suite.router = gin.New()
	v1 := suite.router.Group("/api/v1")
	{
		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
		v1.GET("/users", userHandler.ListUsers)
	}
}

func (suite *UserHandlerTestSuite) TestCreateUser() {
	reqBody := map[string]any{
		"username":  "testuser",
		"email":     "test@example.com",
		"full_name": "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	user := response["user"].(map[string]any)
	assert.Equal(suite.T(), "testuser", user["username"])
	assert.Equal(suite.T(), "test@example.com", user["email"])
	assert.Equal(suite.T(), "Test User", user["full_name"])
	assert.True(suite.T(), user["is_active"].(bool))
	assert.NotEmpty(suite.T(), user["id"])
}

func (suite *UserHandlerTestSuite) TestCreateUserInvalidRequest() {
	reqBody := map[string]any{
		"username": "testuser",
		// Missing email and full_name
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *UserHandlerTestSuite) TestGetUser() {
	// Create a user first
	createReq := &model.CreateUserRequest{
		Username: "getuser",
		Email:    "get@example.com",
		FullName: "Get User",
	}
	user, err := suite.userService.CreateUser(context.Background(), createReq)
	assert.NoError(suite.T(), err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%s", user.ID), nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	responseUser := response["user"].(map[string]any)
	assert.Equal(suite.T(), user.ID, responseUser["id"])
	assert.Equal(suite.T(), user.Username, responseUser["username"])
}

func (suite *UserHandlerTestSuite) TestGetUserNotFound() {
	req, _ := http.NewRequest("GET", "/api/v1/users/nonexistent", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *UserHandlerTestSuite) TestUpdateUser() {
	// Create a user first
	createReq := &model.CreateUserRequest{
		Username: "updateuser",
		Email:    "update@example.com",
		FullName: "Update User",
	}
	user, err := suite.userService.CreateUser(context.Background(), createReq)
	assert.NoError(suite.T(), err)

	reqBody := map[string]any{
		"username":  "updateduser",
		"email":     "updated@example.com",
		"full_name": "Updated User",
		"is_active": false,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%s", user.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	responseUser := response["user"].(map[string]any)
	assert.Equal(suite.T(), "updateduser", responseUser["username"])
	assert.Equal(suite.T(), "updated@example.com", responseUser["email"])
	assert.Equal(suite.T(), "Updated User", responseUser["full_name"])
	assert.False(suite.T(), responseUser["is_active"].(bool))
}

func (suite *UserHandlerTestSuite) TestDeleteUser() {
	// Create a user first
	createReq := &model.CreateUserRequest{
		Username: "deleteuser",
		Email:    "delete@example.com",
		FullName: "Delete User",
	}
	user, err := suite.userService.CreateUser(context.Background(), createReq)
	assert.NoError(suite.T(), err)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%s", user.ID), nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Verify user is deleted
	_, err = suite.userService.GetUser(context.Background(), user.ID)
	assert.Error(suite.T(), err)
}

func (suite *UserHandlerTestSuite) TestListUsers() {
	// Create multiple users
	for i := range 5 {
		createReq := &model.CreateUserRequest{
			Username: fmt.Sprintf("listuser%d", i),
			Email:    fmt.Sprintf("list%d@example.com", i),
			FullName: fmt.Sprintf("List User %d", i),
		}
		_, err := suite.userService.CreateUser(context.Background(), createReq)
		assert.NoError(suite.T(), err)
	}

	req, _ := http.NewRequest("GET", "/api/v1/users?page=1&page_size=3", nil)
	w := httptest.NewRecorder()

	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	users := response["users"].([]any)
	assert.Len(suite.T(), users, 3)
	assert.Equal(suite.T(), float64(5), response["total_count"])
	assert.Equal(suite.T(), float64(1), response["page"])
	assert.Equal(suite.T(), float64(3), response["page_size"])
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}