package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/response"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/validator"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

type authHandler struct {
	s service.AuthService
	v validator.Validator
}

func NewAuthHandler(s service.AuthService, v validator.Validator) *authHandler {
	return &authHandler{s, v}
}

func (h *authHandler) LoginHandler(c *gin.Context) {
	var req service.LoginRequest

	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.Error(apperror.ErrInvalidJSON)
		return
	}

	err = h.v.Validate(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.s.Login(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("login successfully", res))
}
