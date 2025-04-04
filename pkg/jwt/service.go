package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	apperrors "github.com/Brain-Wave-Ecosystem/go-common/pkg/error"
	rand2 "math/rand/v2"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	randMin = 100000
	randMax = 999999
)

type Service struct {
	secret            string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewService(secret string, accessExpiration, refreshExpiration time.Duration) *Service {
	return &Service{
		secret:            secret,
		accessExpiration:  accessExpiration,
		refreshExpiration: refreshExpiration,
	}
}

func (s *Service) GenerateToken(userID string, role string) (string, error) {
	now := time.Now()
	claims := &AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiration)),
		},
		UserID: userID,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *Service) GenerateRefreshToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", fmt.Errorf("failed to generate refresh token %w", err)
	}

	return base64.URLEncoding.EncodeToString(token), nil
}

func (s *Service) GenerateCode() int {
	return rand2.IntN(randMax-randMin) + randMin
}

func (s *Service) GetAccessExpiration() time.Duration {
	return s.accessExpiration
}

func (s *Service) GetRefreshExpiration() time.Duration {
	return s.refreshExpiration
}

func (s *Service) ValidateToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(clearToken(tokenString), &AccessClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(s.secret), nil
	})

	if err != nil && tokenString == "" {
		return nil, apperrors.Forbidden("access token not provided")
	}

	if err != nil {
		return nil, apperrors.Forbidden("invalid token")
	}

	claims, ok := token.Claims.(*AccessClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.Unauthorized("Access token invalid")
}

func (s *Service) ValidateRoleToken(tokenString string, role string) error {
	token, err := jwt.ParseWithClaims(clearToken(tokenString), &AccessClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(s.secret), nil
	})

	if err != nil && tokenString == "" {
		return apperrors.Forbidden("access token not provided")
	}

	if err != nil {
		return apperrors.Forbidden("invalid token")
	}

	claims, ok := token.Claims.(*AccessClaims)
	if ok && token.Valid {
		if claims.Role == string(Admin) {
			return nil
		}

		if claims.Role != role {
			return apperrors.Forbidden("access token invalid")
		}
	}

	return apperrors.Unauthorized("access token invalid")
}

func (s *Service) ValidateUserIDToken(tokenString string, userID string) error {
	token, err := jwt.ParseWithClaims(clearToken(tokenString), &AccessClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(s.secret), nil
	})

	if err != nil && tokenString == "" {
		return apperrors.Forbidden("access token not provided")
	}

	if err != nil {
		return apperrors.Forbidden("invalid token")
	}

	claims, ok := token.Claims.(*AccessClaims)
	if ok && token.Valid {
		if claims.Role == string(Admin) {
			return nil
		}

		if claims.UserID != userID {
			return apperrors.Forbidden("access token invalid")
		}
	}

	return apperrors.Unauthorized("access token invalid")
}

func clearToken(tokenStr string) string {
	return strings.TrimPrefix(tokenStr, "Bearer ")
}
