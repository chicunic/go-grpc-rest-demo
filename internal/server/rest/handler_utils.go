package rest

import (
	"go-grpc-rest-demo/internal/server/errors"
	"go-grpc-rest-demo/internal/server/model"

	"github.com/gin-gonic/gin"
)

func handleUserError(c *gin.Context, err error) {
	appErr := errors.AsAppError(err)
	c.JSON(appErr.ToHTTPStatus(), model.UserResponse{
		Success: false,
		Message: appErr.Message,
	})
}

func respondUserSuccess(c *gin.Context, statusCode int, user *model.User) {
	c.JSON(statusCode, model.UserResponse{
		Success: true,
		User:    user,
		Message: "Operation successful",
	})
}

func handleProductError(c *gin.Context, err error) {
	appErr := errors.AsAppError(err)
	c.JSON(appErr.ToHTTPStatus(), model.ProductResponse{
		Message: appErr.Message,
	})
}

func respondProductSuccess(c *gin.Context, statusCode int, product *model.Product) {
	c.JSON(statusCode, model.ProductResponse{
		Product: product,
		Message: "Operation successful",
	})
}
