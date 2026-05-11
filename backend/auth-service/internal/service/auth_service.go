package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/deliverx/auth-service/internal/middleware"
	"github.com/deliverx/auth-service/internal/model"
	"github.com/deliverx/auth-service/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTokenPrefix = "refresh_token:"
	blacklistPrefix    = "blacklist:"
)

type RegisterRequest struct {
	Name     string     `json:"name" binding:"required"`
	Phone    string     `json:"phone" binding:"required"`
	Email    *string    `json:"email"`
	Password string     `json:"password" binding:"required,min=6"`
	Role     model.Role `json:"role" binding:"required"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role"`
}

type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int64    `json:"expires_in"`
	User         UserInfo `json:"user"`
}

type AuthService struct {
	repo   *repository.UserRepository
	jwt    *middleware.JWTUtil
	redis  *redis.Client
}

func NewAuthService(repo *repository.UserRepository, jwt *middleware.JWTUtil, redis *redis.Client) *AuthService {
	return &AuthService{repo: repo, jwt: jwt, redis: redis}
}

func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	if s.repo.ExistsByPhone(req.Phone) {
		return nil, errors.New("phone number already registered")
	}
	if req.Email != nil && s.repo.ExistsByEmail(*req.Email) {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
		IsActive:     true,
	}
	if err := s.repo.Save(user); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(user)
}

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByPhone(req.Phone)
	if err != nil {
		return nil, errors.New("invalid phone or password")
	}
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid phone or password")
	}
	return s.buildAuthResponse(user)
}

func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	ctx := context.Background()
	key := refreshTokenPrefix + refreshToken
	userIDStr, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	s.redis.Del(ctx, key)
	return s.buildAuthResponse(user)
}

func (s *AuthService) Logout(accessToken, refreshToken string) {
	ctx := context.Background()
	if s.jwt.ValidateToken(accessToken) {
		ttl := time.Duration(s.jwt.ExpirationMs()) * time.Millisecond
		s.redis.Set(ctx, blacklistPrefix+accessToken, "1", ttl)
	}
	if refreshToken != "" {
		s.redis.Del(ctx, refreshTokenPrefix+refreshToken)
	}
}

func (s *AuthService) buildAuthResponse(user *model.User) (*AuthResponse, error) {
	accessToken, err := s.jwt.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	ctx := context.Background()
	ttl := time.Duration(s.jwt.RefreshExpirationMs()) * time.Millisecond
	s.redis.Set(ctx, refreshTokenPrefix+refreshToken, user.ID.String(), ttl)

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwt.ExpirationMs() / 1000,
		User: UserInfo{
			ID:    user.ID.String(),
			Name:  user.Name,
			Phone: user.Phone,
			Email: email,
			Role:  string(user.Role),
		},
	}, nil
}
