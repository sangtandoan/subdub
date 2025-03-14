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
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
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

type (
	CreateUserRequest struct {
		Email    string `json:"email,omitempty"    validate:"email"`
		Password string `json:"password,omitempty" validate:"min=3,max=20"`
	}

	CreateUserResponse struct {
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id,omitempty"`
		Email     string    `json:"email,omitempty"`
		Password  string    `json:"password"`
	}
)

func (s *authService) CreateUser(
	ctx context.Context,
	req *CreateUserRequest,
) (*CreateUserResponse, error) {
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

	return &CreateUserResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		Password:  createdUser.Password,
	}, nil
}
