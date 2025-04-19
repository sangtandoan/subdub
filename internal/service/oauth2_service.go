package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
	"golang.org/x/oauth2"
)

type OAuth2Service interface {
	GenerateURL(ctx context.Context) string
	Callback(ctx context.Context, req *CallBackRequest) (*LoginResponse, error)
}

type googleOAuth2Service struct {
	config           *oauth2.Config
	authProviderRepo repo.AuthProviderRepo
	userRepo         repo.UserRepo
	sessionRepo      repo.SessionRepo
	authenticator    authenticator.Authenticator
	tx               repo.TransactionManager
	state            string
}

func NewGoogleOAuth2Service(
	config *oauth2.Config,
	userRepo repo.UserRepo,
	authProviderRepo repo.AuthProviderRepo,
	sessionRepo repo.SessionRepo,
	authenticator authenticator.Authenticator,
	tx repo.TransactionManager,
) *googleOAuth2Service {
	state, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	return &googleOAuth2Service{
		config,
		authProviderRepo,
		userRepo,
		sessionRepo,
		authenticator,
		tx,
		state.String(),
	}
}

func (s *googleOAuth2Service) GenerateURL(ctx context.Context) string {
	return s.config.AuthCodeURL(s.state, oauth2.AccessTypeOffline)
}

type CallBackRequest struct {
	State string
	Code  string
}

type UserInfoResponse struct {
	ID      string `json:"id,omitempty"`
	Email   string `json:"email,omitempty"`
	Picture string `json:"picture,omitempty"`
}

func (s *googleOAuth2Service) Callback(
	ctx context.Context,
	req *CallBackRequest,
) (*LoginResponse, error) {
	if req.State != s.state {
		return nil, apperror.NewAppError(http.StatusBadRequest, "invalid state")
	}

	token, err := s.config.Exchange(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	res, err := http.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken,
	)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var userInfo UserInfoResponse
	err = json.NewDecoder(res.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	// checks if email exists, if yes -> link auth_provider entry with user
	// if no -> create new user and link with auth_provider entry

	existedUser, err := s.userRepo.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	if err != nil {
		// user does not exists
		err = s.tx.WithTx(ctx, func(txContext context.Context) error {
			// create user
			userID, err := uuid.NewUUID()
			if err != nil {
				return err
			}

			existedUser, err = s.userRepo.CreateUser(txContext, &repo.CreateUserParams{
				ID:    userID,
				Email: userInfo.Email,
			})
			if err != nil {
				return err
			}

			// create auth_provider
			authProviderID, err := uuid.NewUUID()
			if err != nil {
				return err
			}

			_, err = s.authProviderRepo.CreateAuthProvider(
				txContext,
				&repo.CreateAuthProviderParams{
					ID:         authProviderID,
					UserID:     userID,
					Provider:   "google",
					ProviderID: userInfo.ID,
				},
			)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		// user exists -> link auth_provider
		authProviderID, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		_, err = s.authProviderRepo.CreateAuthProvider(
			ctx,
			&repo.CreateAuthProviderParams{
				ID:         authProviderID,
				UserID:     existedUser.ID,
				Provider:   "google",
				ProviderID: userInfo.ID,
			},
		)
		if err != nil {
			return nil, err
		}

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
