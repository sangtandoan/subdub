package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
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

	c.SetCookie(
		"refresh_token",
		res.RefreshToken,
		int(authenticator.RefreshTokenExpiry),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, response.NewAppResponse("login successfully", res))
}

func (h *authHandler) RegisterHandler(c *gin.Context) {
	var req service.RegisterRequest

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

	res, err := h.s.Register(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("created user successfylly", res))
}

func (h *authHandler) LogoutHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.s.Logout(c.Request.Context(), refreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, response.NewAppResponse("logout ok", nil))
}

func (h *authHandler) TokenRenewHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.s.TokenRenew(c.Request.Context(), refreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("renew access token ok", res))
}
