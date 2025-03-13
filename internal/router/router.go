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

	api := g.Group("/api")
	{
		api.Use(middlewares.ErrorMiddleware)

		v1 := api.Group("/v1")
		r.setupAuthRoutes(v1)

		v1.Use(middlewares.AuthMiddleware(r.auth))
		r.setupUserRoutes(v1)
		r.setupSubscriptionRoutes(v1)
	}

	return g
}

func (r *router) setupUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users")

	users.GET("/:id", r.handler.User.GetUserHandler)
	users.POST("", r.handler.User.CreateUserHandler)
}

func (r *router) setupSubscriptionRoutes(group *gin.RouterGroup) {
	sub := group.Group("/subscriptions")

	sub.POST("", r.handler.Subscription.CreateSubscriptionHandler)
	sub.GET(
		"",
		r.handler.Subscription.GetAllSubscriptionsHandler,
	)
}

func (r *router) setupAuthRoutes(group *gin.RouterGroup) {
	auth := group.Group("/auth")

	auth.POST("/login", r.handler.Auth.LoginHandler)
}
