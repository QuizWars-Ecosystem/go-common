package jwt

import (
	"context"
	apperrors "github.com/Brain-Wave-Ecosystem/go-common/pkg/error"
	"google.golang.org/grpc/metadata"
)

const (
	AuthorizationHeader = "authorization"
)

func (s *Service) GenerateTokenWithContext(ctx context.Context, userID, role string) (context.Context, error) {
	token, err := s.GenerateToken(userID, role)
	if err != nil {
		return ctx, err
	}

	return metadata.AppendToOutgoingContext(ctx, AuthorizationHeader, token), nil
}

func (s *Service) ValidateTokenWithContext(ctx context.Context) (*AccessClaims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, apperrors.Unauthorized(AuthAccessTokenNotProvidedError)
	}

	strings := md.Get(AuthorizationHeader)
	if len(strings) == 0 {
		return nil, apperrors.Unauthorized(AuthAccessTokenNotProvidedError)
	}

	claims, err := s.ValidateToken(strings[0])
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *Service) ValidateUserIDWithContext(ctx context.Context, userID string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return apperrors.Unauthorized(AuthAccessTokenNotProvidedError)
	}

	strings := md.Get(AuthorizationHeader)
	if len(strings) == 0 {
		return apperrors.Unauthorized(AuthAccessTokenNotProvidedError)
	}

	return s.ValidateUserIDToken(strings[0], userID)
}

func (s *Service) ValidateRoleWithContext(ctx context.Context, role string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return apperrors.Unauthorized(AuthAccessTokenNotProvidedError)
	}

	strings := md.Get(AuthorizationHeader)
	if len(strings) == 0 {
		return apperrors.Unauthorized(AuthAccessTokenNotProvidedError)
	}

	return s.ValidateRoleToken(strings[0], role)
}
