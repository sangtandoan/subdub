package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/response"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

type userHandler struct {
	s service.UserService
}

func NewUserHandler(s service.UserService) *userHandler {
	return &userHandler{s}
}

func (h *userHandler) GetUserHandler(c *gin.Context) {
	userIDString := c.Param("id")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		_ = c.Error(apperror.ErrInvalidJSON)
		return
	}

	// get userid from context
	userIDFromTokenString, exists := c.Get(authenticator.SubClaim)
	if !exists {
		_ = c.Error(apperror.ErrUnAuthorized)
		return
	}

	userIDFromToken, err := uuid.Parse(userIDFromTokenString.(string))
	if err != nil {
		_ = c.Error(apperror.ErrInvalidUUID)
		return
	}

	// check if authorize with correct userid
	if userID != userIDFromToken {
		_ = c.Error(apperror.ErrUnAuthorized)
		return
	}

	res, err := h.s.GetUser(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(apperror.ErrInvalidJSON)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("get user successfully", res))
}
