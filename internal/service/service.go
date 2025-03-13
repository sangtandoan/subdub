package service

import (
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
)

type Service struct {
	User         UserService
	Subscription SubscriptionService
	Auth         AuthService
}

func NewService(repo *repo.Repo, authenticator authenticator.Authenticator) *Service {
	return &Service{
		User:         NewUserService(repo.User),
		Subscription: NewSubscriptionService(repo.Subscription),
		Auth:         NewAuthService(repo.User, authenticator),
	}
}
