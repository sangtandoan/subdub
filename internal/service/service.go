package service

import (
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/config"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
)

type Service struct {
	User         UserService
	Subscription SubscriptionService
	Auth         AuthService
	OAuth2       OAuth2Service
}

func NewService(
	repo *repo.Repo,
	authenticator authenticator.Authenticator,
	config *config.Config,
) *Service {
	return &Service{
		User:         NewUserService(repo.User),
		Subscription: NewSubscriptionService(repo.Subscription),
		Auth:         NewAuthService(repo.User, authenticator),
		OAuth2:       NewGoogleOAuth2Service(config.GoogleOAuth),
	}
}
