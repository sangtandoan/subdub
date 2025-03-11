package handler

import "github.com/sangtandoan/subscription_tracker/internal/service"

type Handler struct {
	User *userHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		User: NewUserHandler(service.User),
	}
}
