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

//	@BasePath	/api/v1

// GetUserHandler godoc
//
//	@Summary		Get user by id
//	@Description	Get user by id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	service.GetUserResponse
//	@Failure		400	{object}	apperror.AppError
//	@Router			/users/{id} [get]
//	@Security		ApiKeyAuth
func (h *userHandler) GetUserHandler(c *gin.Context) {
	userIDString := c.Param("id")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		_ = c.Error(apperror.ErrInvalidJSON)
		return
	}

	// get userid from context
	userIDFromTokenString, ok := c.Get(authenticator.SubClaim)
	if !ok {
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
