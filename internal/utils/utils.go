package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
)

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	var arr [16]byte

	userIDString, exits := c.Get(authenticator.SubClaim)
	if !exits {
		return uuid.UUID(arr), apperror.ErrUnAuthorized
	}

	return uuid.Parse(userIDString.(string))
}
