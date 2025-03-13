package authenticator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sangtandoan/subscription_tracker/internal/config"
	"github.com/sangtandoan/subscription_tracker/internal/models"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
)

const (
	SubClaim   = "sub"
	EmailClaim = "email"
	ExpClaim   = "exp"
)

// JWT implementation of Authenticator interface
type jwtAuthenticator struct {
	secretKey   string
	tokenExpiry time.Duration
}

func NewJWTAuthenticator(cfg *config.AuthenticatorConfig) (*jwtAuthenticator, error) {
	tokenExpiry, err := time.ParseDuration(cfg.TokenExpiry)
	if err != nil {
		return nil, err
	}

	return &jwtAuthenticator{
		secretKey:   cfg.SecretKey,
		tokenExpiry: tokenExpiry,
	}, nil
}

func (auth *jwtAuthenticator) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		SubClaim:   user.ID,
		EmailClaim: user.Email,
		// exp need convert to unix
		ExpClaim: time.Now().Add(auth.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(auth.secretKey))
}

func (auth *jwtAuthenticator) VerifyToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		// validate signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.NewAppError(
				http.StatusUnauthorized,
				fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]),
			)
		}
		return []byte(auth.secretKey), nil
	},
		jwt.WithExpirationRequired(),
	)
}
