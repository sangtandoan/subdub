package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}

type authService struct {
	userRepo      repo.UserRepo
	authenticator authenticator.Authenticator
}

func NewAuthService(
	userRepo repo.UserRepo,
	authenticator authenticator.Authenticator,
) *authService {
	return &authService{userRepo, authenticator}
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	CreatedAt time.Time
	Token     string
	UserID    uuid.UUID
	Email     string
}

func (s *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// checks user exists?
	existedUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existedUser == nil {
		return nil, apperror.ErrUnAuthorized
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(req.Password))
	if err != nil {
		return nil, apperror.ErrUnAuthorized
	}

	token, err := s.authenticator.GenerateToken(existedUser)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Email:     existedUser.Email,
		UserID:    existedUser.ID,
		CreatedAt: existedUser.CreatedAt,
		Token:     token,
	}, nil
}
