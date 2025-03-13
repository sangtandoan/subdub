package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/response"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/validator"
	"github.com/sangtandoan/subscription_tracker/internal/service"
)

type subscriptionHandler struct {
	s service.SubscriptionService
	v validator.Validator
}

func NewSubscriptionHandler(
	s service.SubscriptionService,
	v validator.Validator,
) *subscriptionHandler {
	return &subscriptionHandler{
		s,
		v,
	}
}

func (h *subscriptionHandler) GetAllSubscriptionsHandler(c *gin.Context) {
	res, err := h.s.GetAllSubscriptions(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("getl all subscriptions successfully", res))
}

func (h *subscriptionHandler) CreateSubscriptionHandler(c *gin.Context) {
	var req service.CreateSubscriptionRequest

	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.v.Validate(req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// userID, exists := c.Get("userID")
	// if !exists {
	// 	_ = c.Error(apperror.NewAppError(http.StatusUnauthorized, "can not find userID in context"))
	// 	return
	// }

	// req.UserID = userID.(uuid.UUID)
	req.UserID, err = uuid.Parse("97182b33-fe85-11ef-bc3e-902e1685779a")
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.s.CreateSubscription(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.NewAppResponse("created subscription sucessfully", res))
}
