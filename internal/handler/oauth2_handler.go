package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/response"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

type oAuth2Handler struct {
	service service.OAuth2Service
}

func NewOAuth2Handler(service service.OAuth2Service) *oAuth2Handler {
	return &oAuth2Handler{service}
}

func (h *oAuth2Handler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"url": "/api/v1/oauth2",
	})
}

func (h *oAuth2Handler) SignInWithOAuth(c *gin.Context) {
	url := h.service.GenerateURL(c.Request.Context())

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *oAuth2Handler) CallbackHandler(c *gin.Context) {
	state := c.Query("state")
	if state == "" {
		c.Error(apperror.NewAppError(http.StatusBadRequest, "invalid callback"))
		return
	}

	code := c.Query("code")
	if state == "" {
		c.Error(apperror.NewAppError(http.StatusBadRequest, "invalid callback"))
		return
	}

	req := &service.CallBackRequest{
		State: state,
		Code:  code,
	}

	res, err := h.service.Callback(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("sign in with oauth ok", res))
}
