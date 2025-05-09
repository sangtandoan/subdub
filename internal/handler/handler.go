package handler

import (
	"github.com/sangtandoan/subscription_tracker/internal/pkg/validator"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

type Handler struct {
	User         *userHandler
	Subscription *subscriptionHandler
	Auth         *authHandler
	OAuth2       *oAuth2Handler
}

func NewHandler(service *service.Service, validator validator.Validator) *Handler {
	return &Handler{
		User:         NewUserHandler(service.User),
		Subscription: NewSubscriptionHandler(service.Subscription, validator),
		Auth:         NewAuthHandler(service.Auth, validator),
		OAuth2:       NewOAuth2Handler(service.OAuth2),
	}
}
