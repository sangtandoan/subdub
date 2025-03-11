package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

type userHandler struct {
	s service.UserService
}

func NewUserHandler(s service.UserService) *userHandler {
	return &userHandler{s}
}

func (h *userHandler) CreateUserHandler(c *gin.Context) {
	var req service.CreateUserRequest

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "json format invalid"})
	}

	res, err := h.s.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": res})
}

func (h *userHandler) GetUserHandler(c *gin.Context) {
	userIDString := c.Param("id")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		apperror.WriteError(c, err)
	}

	res, err := h.s.GetUser(c.Request.Context(), userID)
	if err != nil {
		apperror.WriteError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": res})
}
