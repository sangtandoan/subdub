package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/response"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/validator"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

type userHandler struct {
	s service.UserService
	v validator.Validator
}

func NewUserHandler(s service.UserService, v validator.Validator) *userHandler {
	return &userHandler{s, v}
}

func (h *userHandler) CreateUserHandler(c *gin.Context) {
	var req service.CreateUserRequest

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

	res, err := h.s.CreateUser(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("created user successfylly", res))
}

func (h *userHandler) GetUserHandler(c *gin.Context) {
	userIDString := c.Param("id")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		_ = c.Error(apperror.ErrInvalidJSON)
		return
	}

	res, err := h.s.GetUser(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(apperror.ErrInvalidJSON)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("get user successfylly", res))
}
