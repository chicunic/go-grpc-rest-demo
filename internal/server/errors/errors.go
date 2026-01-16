package errors

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	// Client errors
	ErrCodeInvalidRequest   ErrorCode = "INVALID_REQUEST"
	ErrCodeValidationFailed ErrorCode = "VALIDATION_FAILED"
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists    ErrorCode = "ALREADY_EXISTS"
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"

	// Server errors
	ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceDown    ErrorCode = "SERVICE_DOWN"
	ErrCodeDatabaseError  ErrorCode = "DATABASE_ERROR"
	ErrCodeExternalAPI    ErrorCode = "EXTERNAL_API_ERROR"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	Field      string    `json:"field,omitempty"`
	HTTPStatus int       `json:"-"`
	GRPCCode   codes.Code `json:"-"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ToHTTPStatus returns the appropriate HTTP status code
func (e *AppError) ToHTTPStatus() int {
	if e.HTTPStatus != 0 {
		return e.HTTPStatus
	}

	switch e.Code {
	case ErrCodeInvalidRequest, ErrCodeValidationFailed:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeAlreadyExists:
		return http.StatusConflict
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeInternal, ErrCodeDatabaseError:
		return http.StatusInternalServerError
	case ErrCodeServiceDown, ErrCodeExternalAPI:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// ToGRPCStatus returns the appropriate gRPC status
func (e *AppError) ToGRPCStatus() *status.Status {
	if e.GRPCCode != codes.OK {
		return status.New(e.GRPCCode, e.Message)
	}

	var grpcCode codes.Code
	switch e.Code {
	case ErrCodeInvalidRequest, ErrCodeValidationFailed:
		grpcCode = codes.InvalidArgument
	case ErrCodeNotFound:
		grpcCode = codes.NotFound
	case ErrCodeAlreadyExists:
		grpcCode = codes.AlreadyExists
	case ErrCodeUnauthorized:
		grpcCode = codes.Unauthenticated
	case ErrCodeForbidden:
		grpcCode = codes.PermissionDenied
	case ErrCodeInternal, ErrCodeDatabaseError:
		grpcCode = codes.Internal
	case ErrCodeServiceDown, ErrCodeExternalAPI:
		grpcCode = codes.Unavailable
	default:
		grpcCode = codes.Internal
	}

	return status.New(grpcCode, e.Message)
}

// Error constructors for common scenarios

func NewValidationError(field, message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidationFailed,
		Message: message,
		Field:   field,
	}
}

func NewNotFoundError(resource, id string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Details: fmt.Sprintf("ID: %s", id),
	}
}

func NewAlreadyExistsError(resource, field, value string) *AppError {
	return &AppError{
		Code:    ErrCodeAlreadyExists,
		Message: fmt.Sprintf("%s already exists", resource),
		Details: fmt.Sprintf("%s: %s", field, value),
		Field:   field,
	}
}

func NewInvalidRequestError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeInvalidRequest,
		Message: message,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeInternal,
		Message: message,
	}
}

func NewDatabaseError(operation string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeDatabaseError,
		Message: "Database operation failed",
		Details: fmt.Sprintf("Operation: %s, Error: %v", operation, err),
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError converts an error to AppError, or creates a generic internal error
func AsAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	return &AppError{
		Code:    ErrCodeInternal,
		Message: "Internal server error",
		Details: err.Error(),
	}
}