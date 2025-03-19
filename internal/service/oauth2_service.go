package service

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"golang.org/x/oauth2"
)

type OAuth2Service interface {
	GenerateURL(ctx context.Context) string
	Callback(ctx context.Context, req *CallBackRequest) (*CallbackResponse, error)
}

type googleOAuth2Service struct {
	config *oauth2.Config
	state  string
}

func NewGoogleOAuth2Service(config *oauth2.Config) *googleOAuth2Service {
	state, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	return &googleOAuth2Service{config, state.String()}
}

func (s *googleOAuth2Service) GenerateURL(ctx context.Context) string {
	return s.config.AuthCodeURL(s.state, oauth2.AccessTypeOffline)
}

type CallBackRequest struct {
	State string
	Code  string
}

type CallbackResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (s *googleOAuth2Service) Callback(
	ctx context.Context,
	req *CallBackRequest,
) (*CallbackResponse, error) {
	if req.State != s.state {
		return nil, apperror.NewAppError(http.StatusBadRequest, "invalid state")
	}

	token, err := s.config.Exchange(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	return &CallbackResponse{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken}, nil
}
