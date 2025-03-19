package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/handler"
	"github.com/sangtandoan/subscription_tracker/internal/middlewares"
)

type router struct {
	handler *handler.Handler
	auth    authenticator.Authenticator
}

func NewRouter(handler *handler.Handler, auth authenticator.Authenticator) *router {
	return &router{handler, auth}
}

func (r *router) Setup() http.Handler {
	g := gin.Default()
	g.LoadHTMLGlob("templates/*")

	api := g.Group("/api")
	{
		api.Use(middlewares.ErrorMiddleware)
		api.Use(middlewares.GZipMiddleware)

		v1 := api.Group("/v1")
		{
			r.setupOAuthRoutes(v1)
			r.setupAuthRoutes(v1)

			// protected routes
			v1.Use(middlewares.AuthMiddleware(r.auth))
			r.setupUserRoutes(v1)
			r.setupSubscriptionRoutes(v1)
		}
	}

	return g
}

func (r *router) setupUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users")

	users.GET("/:id", r.handler.User.GetUserHandler)
}

func (r *router) setupSubscriptionRoutes(group *gin.RouterGroup) {
	sub := group.Group("/subscriptions")

	sub.POST("", r.handler.Subscription.CreateSubscriptionHandler)
	sub.GET("", r.handler.Subscription.GetAllSubscriptionsHandler)
	// sub.GET("", r.handler.Subscription.GetSubscriptionsBeforeNumDays)
}

func (r *router) setupAuthRoutes(group *gin.RouterGroup) {
	auth := group.Group("/auth")

	auth.POST("/login", r.handler.Auth.LoginHandler)
	auth.POST("/register", r.handler.Auth.RegisterHandler)
	auth.POST("/logout", r.handler.Auth.LogoutHandler)
	auth.POST("/renew", r.handler.Auth.TokenRenewHandler)
}

func (r *router) setupOAuthRoutes(group *gin.RouterGroup) {
	oauth2 := group.Group("/oauth2")

	oauth2.GET("", r.handler.OAuth2.SignInWithOAuth)
	oauth2.GET("callback", r.handler.OAuth2.CallbackHandler)
	oauth2.GET("login", r.handler.OAuth2.Login)
}
