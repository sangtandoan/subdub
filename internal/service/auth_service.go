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
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
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
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token,omitempty"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Email     string    `json:"email,omitempty"`
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

type (
	RegisterRequest struct {
		Email    string `json:"email,omitempty"    validate:"email"`
		Password string `json:"password,omitempty" validate:"min=3,max=20"`
	}

	RegisterResponse struct {
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id,omitempty"`
		Email     string    `json:"email,omitempty"`
		Password  string    `json:"password"`
	}
)

func (s *authService) Register(
	ctx context.Context,
	req *RegisterRequest,
) (*RegisterResponse, error) {
	// check if user exists
	existed, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existed != nil {
		return nil, apperror.ErrExisted
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// generate userID
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	createdUser, err := s.userRepo.CreateUser(ctx, &repo.CreateUserParams{
		ID:       userID,
		Email:    req.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		Password:  createdUser.Password,
	}, nil
}
