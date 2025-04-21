package utils

import (
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/models"
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

func RandomUser() *models.User {
	return &models.User{
		CreatedAt: time.Now(),
		Email:     gofakeit.Email(),
		Password:  gofakeit.Password(true, true, true, false, false, 10),
		ID:        uuid.New(),
	}
}
