package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	rand2 "math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/QuizWars-Ecosystem/go-common/pkg/abstractions"

	apperrors "github.com/Brain-Wave-Ecosystem/go-common/pkg/error"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	randMin = 100000
	randMax = 999999
)

const (
	AuthAccessTokenNotProvidedError = "access token not provided"
	AuthInvalidTokenError           = "invalid token"
	AuthPermissionDeniedError       = "permission denied"
	AuthAccessTokenInvalid          = "access token invalid"
)

var ultimativeRoleList = []string{string(Admin), string(Super)}

var _ abstractions.ConfigSubscriber[*Config] = (*Service)(nil)

type Config struct {
	Secret            string        `mapstructure:"Secret"`
	AccessExpiration  time.Duration `mapstructure:"access_expiration"`
	RefreshExpiration time.Duration `mapstructure:"refresh_expiration"`
}

type Service struct {
	cfg *Config
	mx  sync.RWMutex
}

func NewService(cfg *Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) SectionKey() string {
	return "JWT"
}

func (s *Service) UpdateConfig(newCfg *Config) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.cfg = newCfg

	return nil
}

func (s *Service) GenerateToken(userID string, role string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	now := time.Now()
	claims := &AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessExpiration)),
		},
		UserID: userID,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
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
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.cfg.AccessExpiration
}

func (s *Service) GetRefreshExpiration() time.Duration {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.cfg.RefreshExpiration
}

func (s *Service) ValidateToken(tokenString string) (*AccessClaims, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	token, err := jwt.ParseWithClaims(clearToken(tokenString), &AccessClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(s.cfg.Secret), nil
	})

	if err != nil && tokenString == "" {
		return nil, apperrors.Forbidden(AuthAccessTokenNotProvidedError)
	}

	if err != nil {
		return nil, apperrors.Forbidden(AuthInvalidTokenError)
	}

	claims, ok := token.Claims.(*AccessClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.Unauthorized(AuthAccessTokenInvalid)
}

func (s *Service) ValidateRoleToken(tokenString string, role string) error {
	s.mx.RLock()
	defer s.mx.RUnlock()
	token, err := jwt.ParseWithClaims(clearToken(tokenString), &AccessClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(s.cfg.Secret), nil
	})

	if err != nil && tokenString == "" {
		return apperrors.Forbidden(AuthAccessTokenNotProvidedError)
	}

	if err != nil {
		return apperrors.Forbidden(AuthInvalidTokenError)
	}

	claims, ok := token.Claims.(*AccessClaims)
	if ok && token.Valid {
		if ultimativeRoles(claims.Role) {
			return nil
		}

		if claims.Role != role {
			return apperrors.Forbidden(AuthPermissionDeniedError)
		}

		return nil
	}

	return apperrors.Unauthorized(AuthAccessTokenInvalid)
}

func (s *Service) ValidateUserIDToken(tokenString string, userID string) error {
	s.mx.RLock()
	defer s.mx.RUnlock()
	token, err := jwt.ParseWithClaims(clearToken(tokenString), &AccessClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(s.cfg.Secret), nil
	})

	if err != nil && tokenString == "" {
		return apperrors.Forbidden(AuthAccessTokenNotProvidedError)
	}

	if err != nil {
		return apperrors.Forbidden(AuthInvalidTokenError)
	}

	claims, ok := token.Claims.(*AccessClaims)
	if ok && token.Valid {
		if ultimativeRoles(claims.Role) {
			return nil
		}

		if claims.UserID != userID {
			return apperrors.Forbidden(AuthPermissionDeniedError)
		}

		return nil
	}

	return apperrors.Unauthorized(AuthAccessTokenInvalid)
}

func clearToken(tokenStr string) string {
	return strings.TrimPrefix(tokenStr, "Bearer ")
}

func ultimativeRoles(role string) bool {
	for _, r := range ultimativeRoleList {
		if role == r {
			return true
		}
	}

	return false
}
