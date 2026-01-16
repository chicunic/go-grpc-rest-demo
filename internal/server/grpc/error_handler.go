package grpc

import (
	"go-grpc-rest-demo/internal/server/errors"
)

// handleGRPCError converts AppError to gRPC status error
func handleGRPCError(err error) error {
	if err == nil {
		return nil
	}
	return errors.AsAppError(err).ToGRPCStatus().Err()
}