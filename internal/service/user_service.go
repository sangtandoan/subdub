package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*GetUserResponse, error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
}

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) *userService {
	return &userService{
		userRepo,
	}
}

type GetUserResponse struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password"`
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*GetUserResponse, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}, nil
}

type (
	CreateUserRequest struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}

	CreateUserResponse struct {
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id,omitempty"`
		Email     string    `json:"email,omitempty"`
		Password  string    `json:"password"`
	}
)

func (s *userService) CreateUser(
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
		return nil, err
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
