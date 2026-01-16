package grpc

import (
	"context"
	"fmt"
	"testing"

	pb "go-grpc-rest-demo/api/gen/go/user/v1"
	"go-grpc-rest-demo/internal/server/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserServerTestSuite struct {
	suite.Suite
	server *UserServer
}

func (suite *UserServerTestSuite) SetupTest() {
	userService := service.NewUserService()
	suite.server = NewUserServer(userService)
}

func (suite *UserServerTestSuite) TestCreateUser() {
	req := &pb.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		FullName: "Test User",
	}

	resp, err := suite.server.CreateUser(context.Background(), req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.NotNil(suite.T(), resp.User)
	assert.Equal(suite.T(), req.Username, resp.User.Username)
	assert.Equal(suite.T(), req.Email, resp.User.Email)
	assert.Equal(suite.T(), req.FullName, resp.User.FullName)
	assert.True(suite.T(), resp.User.IsActive)
	assert.NotEmpty(suite.T(), resp.User.Id)
	assert.NotEmpty(suite.T(), resp.User.CreatedAt)
	assert.NotEmpty(suite.T(), resp.User.UpdatedAt)
}

func (suite *UserServerTestSuite) TestCreateUserDuplicate() {
	req := &pb.CreateUserRequest{
		Username: "duplicate",
		Email:    "duplicate@example.com",
		FullName: "Duplicate User",
	}

	_, err := suite.server.CreateUser(context.Background(), req)
	assert.NoError(suite.T(), err)

	_, err = suite.server.CreateUser(context.Background(), req)
	assert.Error(suite.T(), err)
}

func (suite *UserServerTestSuite) TestGetUser() {
	// Create a user first
	createReq := &pb.CreateUserRequest{
		Username: "getuser",
		Email:    "get@example.com",
		FullName: "Get User",
	}

	createResp, err := suite.server.CreateUser(context.Background(), createReq)
	assert.NoError(suite.T(), err)

	// Get the user
	getReq := &pb.GetUserRequest{
		Id: createResp.User.Id,
	}

	getResp, err := suite.server.GetUser(context.Background(), getReq)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), getResp)
	assert.NotNil(suite.T(), getResp.User)
	assert.Equal(suite.T(), createResp.User.Id, getResp.User.Id)
	assert.Equal(suite.T(), createResp.User.Username, getResp.User.Username)
}

func (suite *UserServerTestSuite) TestGetUserNotFound() {
	req := &pb.GetUserRequest{
		Id: "nonexistent",
	}

	_, err := suite.server.GetUser(context.Background(), req)
	assert.Error(suite.T(), err)
}

func (suite *UserServerTestSuite) TestUpdateUser() {
	// Create a user first
	createReq := &pb.CreateUserRequest{
		Username: "updateuser",
		Email:    "update@example.com",
		FullName: "Update User",
	}

	createResp, err := suite.server.CreateUser(context.Background(), createReq)
	assert.NoError(suite.T(), err)

	// Update the user
	newUsername := "updateduser"
	newEmail := "updated@example.com"
	newFullName := "Updated User"
	isActive := false

	updateReq := &pb.UpdateUserRequest{
		Id:       createResp.User.Id,
		Username: &newUsername,
		Email:    &newEmail,
		FullName: &newFullName,
		IsActive: &isActive,
	}

	updateResp, err := suite.server.UpdateUser(context.Background(), updateReq)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updateResp)
	assert.NotNil(suite.T(), updateResp.User)
	assert.Equal(suite.T(), newUsername, updateResp.User.Username)
	assert.Equal(suite.T(), newEmail, updateResp.User.Email)
	assert.Equal(suite.T(), newFullName, updateResp.User.FullName)
	assert.Equal(suite.T(), isActive, updateResp.User.IsActive)
}

func (suite *UserServerTestSuite) TestDeleteUser() {
	// Create a user first
	createReq := &pb.CreateUserRequest{
		Username: "deleteuser",
		Email:    "delete@example.com",
		FullName: "Delete User",
	}

	createResp, err := suite.server.CreateUser(context.Background(), createReq)
	assert.NoError(suite.T(), err)

	// Delete the user
	deleteReq := &pb.DeleteUserRequest{
		Id: createResp.User.Id,
	}

	_, err = suite.server.DeleteUser(context.Background(), deleteReq)
	assert.NoError(suite.T(), err)

	// Verify user is deleted
	getReq := &pb.GetUserRequest{
		Id: createResp.User.Id,
	}

	_, err = suite.server.GetUser(context.Background(), getReq)
	assert.Error(suite.T(), err)
}

func (suite *UserServerTestSuite) TestListUsers() {
	// Create multiple users
	usernames := []string{"list1", "list2", "list3", "list4", "list5"}
	for i, username := range usernames {
		createReq := &pb.CreateUserRequest{
			Username: username,
			Email:    fmt.Sprintf("list%d@example.com", i+1),
			FullName: fmt.Sprintf("List User %d", i+1),
		}
		_, err := suite.server.CreateUser(context.Background(), createReq)
		assert.NoError(suite.T(), err)
	}

	// Test pagination
	listReq := &pb.ListUsersRequest{
		Page:     1,
		PageSize: 3,
	}

	listResp, err := suite.server.ListUsers(context.Background(), listReq)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), listResp)
	assert.Len(suite.T(), listResp.Users, 3)
	assert.Equal(suite.T(), int32(5), listResp.TotalCount)
	assert.Equal(suite.T(), int32(1), listResp.Page)
	assert.Equal(suite.T(), int32(3), listResp.PageSize)
}

func TestUserServerTestSuite(t *testing.T) {
	suite.Run(t, new(UserServerTestSuite))
}