package authenticator

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sangtandoan/subscription_tracker/internal/models"
)

const (
	RefreshTokenExpiry = time.Hour * 24 * 30 // 30 days
)

type Authenticator interface {
	GenerateToken(user *models.User) (string, error)
	GenerateRefreshToken(user *models.User, sessionID string) (string, time.Time, error)
	VerifyToken(token string) (*jwt.Token, error)
}
