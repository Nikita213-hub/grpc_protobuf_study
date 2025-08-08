package auth

import (
	"errors"

	authservice "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/services/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidEmail   = errors.New("invalid email format")
	ErrEmptyEmail     = errors.New("email is required")
	ErrEmptyCode      = errors.New("verification code is required")
	ErrInvalidCodeLen = errors.New("verification code must be 6 digits")
)

func convertServiceError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, authservice.ErrIncorrectCode):
		return status.Error(codes.Unauthenticated, "verification code is incorrect")
	case errors.Is(err, authservice.ErrCodeExpired):
		return status.Error(codes.Unauthenticated, "verification code has expired")
	case errors.Is(err, authservice.ErrSessionExpired):
		return status.Error(codes.Unauthenticated, "session has expired")
	case errors.Is(err, authservice.ErrSessionNotFound):
		return status.Error(codes.NotFound, "session not found")

	case errors.Is(err, authservice.ErrCodeSaveFailed):
		return status.Error(codes.Internal, "service temporarily unavailable")
	case errors.Is(err, authservice.ErrSessionCreateFailed):
		return status.Error(codes.Internal, "service temporarily unavailable")

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
