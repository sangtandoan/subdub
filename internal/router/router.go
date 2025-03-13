package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sangtandoan/subscription_tracker/internal/handler"
	"github.com/sangtandoan/subscription_tracker/internal/middlewares"
)

type router struct {
	handler *handler.Handler
}

func NewRouter(handler *handler.Handler) *router {
	return &router{handler}
}

func (r *router) Setup() http.Handler {
	g := gin.Default()

	api := g.Group("/api")
	{
		api.Use(middlewares.ErrorMiddleware)

		v1 := api.Group("/v1")
		r.setupUserRoutes(v1)
		r.setupSubscriptionRoutes(v1)
	}

	return g
}

func (r *router) setupUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users")

	// users.GET("", r.handler.User.GetUserHandler)
	users.GET("/:id", r.handler.User.GetUserHandler)
	users.POST("", r.handler.User.CreateUserHandler)
}

func (r *router) setupSubscriptionRoutes(group *gin.RouterGroup) {
	sub := group.Group("/subscriptions")

	sub.POST("", r.handler.Subscription.CreateSubscriptionHandler)
	sub.GET("", r.handler.Subscription.GetAllSubscriptionsHandler)
}
