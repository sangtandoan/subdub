package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
)

type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*GetUserResponse, error)
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
