package handler

import (
	"github.com/sangtandoan/subscription_tracker/internal/pkg/validator"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

type Handler struct {
	User         *userHandler
	Subscription *subscriptionHandler
	Auth         *authHandler
}

func NewHandler(service *service.Service, validator validator.Validator) *Handler {
	return &Handler{
		User:         NewUserHandler(service.User, validator),
		Subscription: NewSubscriptionHandler(service.Subscription, validator),
		Auth:         NewAuthHandler(service.Auth, validator),
	}
}
