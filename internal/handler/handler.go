package handler

import (
	"github.com/sangtandoan/subscription_tracker/internal/pkg/validator"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

type Handler struct {
	User *userHandler
}

func NewHandler(service *service.Service, validator validator.Validator) *Handler {
	return &Handler{
		User: NewUserHandler(service.User, validator),
	}
}
