package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/models"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	TokenRenew(ctx context.Context, refreshToken string) (*TokenRenewResponse, error)
}

type authService struct {
	userRepo      repo.UserRepo
	sessionRepo   repo.SessionRepo
	authenticator authenticator.Authenticator
}

func NewAuthService(
	userRepo repo.UserRepo,
	sessionRepo repo.SessionRepo,
	authenticator authenticator.Authenticator,
) *authService {
	return &authService{userRepo, sessionRepo, authenticator}
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Token        string    `json:"token,omitempty"`
	Email        string    `json:"email,omitempty"`
	SessionID    string    `json:"-"`
	RefreshToken string    `json:"-"`
	UserID       uuid.UUID `json:"user_id,omitempty"`
}

func (s *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// checks user exists?
	existedUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existedUser == nil {
		return nil, apperror.ErrUnAuthorized
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(req.Password))
	if err != nil {
		return nil, apperror.ErrUnAuthorized
	}

	tokens, err := generateTokens(ctx, s.authenticator, s.sessionRepo, existedUser)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Email:        existedUser.Email,
		UserID:       existedUser.ID,
		Token:        tokens.accessToken,
		RefreshToken: tokens.refreshToken,
		SessionID:    tokens.session.ID.String(),
	}, nil
}

type (
	RegisterRequest struct {
		Email    string `json:"email,omitempty"    validate:"email"`
		Password string `json:"password,omitempty" validate:"min=3,max=20"`
	}

	RegisterResponse struct {
		CreatedAt time.Time `json:"created_at"`
		ID        uuid.UUID `json:"id,omitempty"`
		Email     string    `json:"email,omitempty"`
		Password  string    `json:"password"`
	}
)

func (s *authService) Register(
	ctx context.Context,
	req *RegisterRequest,
) (*RegisterResponse, error) {
	// check if user exists
	existed, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if existed != nil {
		return nil, apperror.ErrExisted
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// generate userID
	userID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	createdUser, err := s.userRepo.CreateUser(ctx, &repo.CreateUserParams{
		ID:       userID,
		Email:    req.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		Password:  createdUser.Password,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.getRefreshTokenClaims(refreshToken)
	if err != nil {
		return err
	}

	sessionID, err := s.getSessionIDFromClaims(claims)
	if err != nil {
		return err
	}

	if err := s.sessionRepo.DeleteSession(ctx, sessionID); err != nil {
		return err
	}
	// when logout, delete session and clear cookie contains refrsh token
	// access token has short lifetime so just wait and it will expire

	return nil
}

type TokenRenewResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *authService) TokenRenew(
	ctx context.Context,
	refreshToken string,
) (*TokenRenewResponse, error) {
	claims, err := s.getRefreshTokenClaims(refreshToken)
	if err != nil {
		return nil, err
	}

	email, err := s.getUserEmailFromClaims(claims)
	if err != nil {
		return nil, err
	}

	sessionID, err := s.getSessionIDFromClaims(claims)
	if err != nil {
		return nil, err
	}

	session, err := s.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// checks if email matches
	if email != session.UserEmail {
		return nil, apperror.ErrUnAuthorized
	}

	// checks if refreshToken matches
	if refreshToken != session.RefreshToken {
		return nil, apperror.ErrUnAuthorized
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.authenticator.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenRenewResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *authService) getRefreshTokenClaims(refreshToken string) (jwt.MapClaims, error) {
	token, err := s.authenticator.VerifyToken(refreshToken)
	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, apperror.ErrUnAuthorized
	}

	return claims, nil
}

func (s *authService) getSessionIDFromClaims(claims jwt.MapClaims) (uuid.UUID, error) {
	sessionIDStr, err := claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return uuid.UUID{}, err
	}

	return sessionID, nil
}

func (s *authService) getUserEmailFromClaims(claims jwt.MapClaims) (string, error) {
	email, ok := claims[authenticator.EmailClaim]
	if !ok {
		return "", apperror.ErrUnAuthorized
	}
	return email.(string), nil
}

type generateTokensResponse struct {
	accessToken  string
	refreshToken string
	session      *models.Session
}

func generateTokens(
	ctx context.Context,
	authenticator authenticator.Authenticator,
	sessionRepo repo.SessionRepo,
	existedUser *models.User,
) (*generateTokensResponse, error) {
	accessToken, err := authenticator.GenerateToken(existedUser)
	if err != nil {
		return nil, err
	}

	sessionID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	refreshToken, expiresAt, err := authenticator.GenerateRefreshToken(
		existedUser,
		sessionID.String(),
	)
	if err != nil {
		return nil, err
	}

	session, err := sessionRepo.CreateSession(ctx, &repo.CreateSessionParams{
		ID:           sessionID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		UserEmail:    existedUser.Email,
	})
	if err != nil {
		return nil, err
	}

	return &generateTokensResponse{
		accessToken,
		refreshToken,
		session,
	}, nil
}
