package errors

import (
	"google.golang.org/grpc/status"
	"testing"

	apperrors "github.com/Brain-Wave-Ecosystem/go-common/pkg/error"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func RequireNotFoundError(t *testing.T, err error, subject, key string, value any) {
	msg := apperrors.NotFound(subject, key, value).Error()
	requireAPIError(t, err, codes.NotFound, msg)
}

func RequireUnauthorizedError(t *testing.T, err error, msg string) {
	requireAPIError(t, err, codes.Unauthenticated, msg)
}

func RequireForbiddenError(t *testing.T, err error, msg string) {
	requireAPIError(t, err, codes.PermissionDenied, msg)
}

func RequireBadRequestError(t *testing.T, err error, msg string) {
	requireAPIError(t, err, codes.InvalidArgument, msg)
}

func RequireAlreadyExistsError(t *testing.T, err error, subject, key string, value any) {
	msg := apperrors.AlreadyExists(subject, key, value).Error()
	requireAPIError(t, err, codes.AlreadyExists, msg)
}

func requireAPIError(t *testing.T, err error, code codes.Code, msg string) {
	s, ok := status.FromError(err)
	require.True(t, ok, "expected grpc status error")
	require.Equal(t, code, s.Code())
	require.Contains(t, s.Message(), msg)
}
