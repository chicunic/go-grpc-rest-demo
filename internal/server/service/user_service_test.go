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

type UserServiceTestSuite struct {
	suite.Suite
	service *UserService
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.service = NewUserService()
}

func (suite *UserServiceTestSuite) TestCreateUser() {
	req := &model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
	}

	user, err := suite.service.CreateUser(context.Background(), req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), req.Username, user.Username)
	assert.Equal(suite.T(), req.Email, user.Email)
	assert.Equal(suite.T(), req.FullName, user.FullName)
	assert.True(suite.T(), user.IsActive)
	assert.NotEmpty(suite.T(), user.ID)
	assert.WithinDuration(suite.T(), time.Now(), user.CreatedAt, time.Second)
	assert.WithinDuration(suite.T(), time.Now(), user.UpdatedAt, time.Second)
}

func (suite *UserServiceTestSuite) TestCreateUserDuplicateUsername() {
	req1 := &model.CreateUserRequest{
		Username: "duplicate",
		Email:    "test1@example.com",
		FullName: "Test User",
	}
	req2 := &model.CreateUserRequest{
		Username: "duplicate",
		Email:    "test2@example.com",
		FullName: "Test User",
	}

	_, err := suite.service.CreateUser(context.Background(), req1)
	assert.NoError(suite.T(), err)

	_, err = suite.service.CreateUser(context.Background(), req2)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "already exists")
}

func (suite *UserServiceTestSuite) TestCreateUserDuplicateEmail() {
	req1 := &model.CreateUserRequest{
		Username: "user1",
		Email:    "duplicate@example.com",
		FullName: "Test User",
	}
	req2 := &model.CreateUserRequest{
		Username: "user2",
		Email:    "duplicate@example.com",
		FullName: "Test User",
	}

	_, err := suite.service.CreateUser(context.Background(), req1)
	assert.NoError(suite.T(), err)

	_, err = suite.service.CreateUser(context.Background(), req2)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "already exists")
}

func (suite *UserServiceTestSuite) TestGetUser() {
	req := &model.CreateUserRequest{
		Username: "getuser",
		Email:    "get@example.com",
		FullName: "Get User",
	}

	created, err := suite.service.CreateUser(context.Background(), req)
	assert.NoError(suite.T(), err)

	retrieved, err := suite.service.GetUser(context.Background(), created.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), created.ID, retrieved.ID)
	assert.Equal(suite.T(), created.Username, retrieved.Username)
	assert.Equal(suite.T(), created.Email, retrieved.Email)
}

func (suite *UserServiceTestSuite) TestGetUserNotFound() {
	_, err := suite.service.GetUser(context.Background(), "nonexistent")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *UserServiceTestSuite) TestUpdateUser() {
	createReq := &model.CreateUserRequest{
		Username: "updateuser",
		Email:    "update@example.com",
		FullName: "Update User",
	}
	user, err := suite.service.CreateUser(context.Background(), createReq)
	assert.NoError(suite.T(), err)

	newUsername := "updateduser"
	newEmail := "updated@example.com"
	newFullName := "Updated User"
	isActive := false

	updateReq := &model.UpdateUserRequest{
		ID:       user.ID,
		Username: &newUsername,
		Email:    &newEmail,
		FullName: &newFullName,
		IsActive: &isActive,
	}

	updated, err := suite.service.UpdateUser(context.Background(), updateReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newUsername, updated.Username)
	assert.Equal(suite.T(), newEmail, updated.Email)
	assert.Equal(suite.T(), newFullName, updated.FullName)
	assert.Equal(suite.T(), isActive, updated.IsActive)
	assert.True(suite.T(), updated.UpdatedAt.After(updated.CreatedAt))
}

func (suite *UserServiceTestSuite) TestDeleteUser() {
	createReq := &model.CreateUserRequest{
		Username: "deleteuser",
		Email:    "delete@example.com",
		FullName: "Delete User",
	}
	user, err := suite.service.CreateUser(context.Background(), createReq)
	assert.NoError(suite.T(), err)

	err = suite.service.DeleteUser(context.Background(), user.ID)
	assert.NoError(suite.T(), err)

	_, err = suite.service.GetUser(context.Background(), user.ID)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")
}

func (suite *UserServiceTestSuite) TestListUsers() {
	for i := 0; i < 5; i++ {
		createReq := &model.CreateUserRequest{
			Username: fmt.Sprintf("listuser%d", i),
			Email:    fmt.Sprintf("list%d@example.com", i),
			FullName: fmt.Sprintf("List User %d", i),
		}
		_, err := suite.service.CreateUser(context.Background(), createReq)
		assert.NoError(suite.T(), err)
	}

	listReq := &model.ListUsersRequest{
		Page:     1,
		PageSize: 3,
	}
	result, totalCount, page, pageSize, err := suite.service.ListUsers(context.Background(), listReq)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)
	assert.Equal(suite.T(), int32(5), totalCount)
	assert.Equal(suite.T(), int32(1), page)
	assert.Equal(suite.T(), int32(3), pageSize)

	listReq.Page = 2
	result, totalCount, page, pageSize, err = suite.service.ListUsers(context.Background(), listReq)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), int32(5), totalCount)
	assert.Equal(suite.T(), int32(2), page)
	assert.Equal(suite.T(), int32(3), pageSize)
}

func (suite *UserServiceTestSuite) TestListUsersWithFilter() {
	users := []model.CreateUserRequest{
		{Username: "alice", Email: "alice@example.com", FullName: "Alice Smith"},
		{Username: "bob", Email: "bob@example.com", FullName: "Bob Jones"},
		{Username: "charlie", Email: "charlie@example.com", FullName: "Charlie Brown"},
	}

	for _, req := range users {
		_, err := suite.service.CreateUser(context.Background(), &req)
		assert.NoError(suite.T(), err)
	}

	filter := "alice"
	listReq := &model.ListUsersRequest{
		Page:     1,
		PageSize: 10,
		Filter:   &filter,
	}
	result, totalCount, _, _, err := suite.service.ListUsers(context.Background(), listReq)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), int32(1), totalCount)
	assert.Equal(suite.T(), "alice", result[0].Username)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}