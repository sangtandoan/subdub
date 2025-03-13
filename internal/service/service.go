package service

import "github.com/sangtandoan/subscription_tracker/internal/repo"

type Service struct {
	User         UserService
	Subscription SubscriptionService
}

func NewService(repo *repo.Repo) *Service {
	return &Service{
		User:         NewUserService(repo.User),
		Subscription: NewSubscriptionService(repo.Subscription),
	}
}
