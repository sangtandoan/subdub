package authenticator

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/sangtandoan/subscription_tracker/internal/models"
)

type Authenticator interface {
	GenerateToken(user *models.User) (string, error)
	VerifyToken(token string) (*jwt.Token, error)
}
