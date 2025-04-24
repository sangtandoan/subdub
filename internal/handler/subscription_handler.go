package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/response"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/validator"
	"github.com/sangtandoan/subscription_tracker/internal/service"
	"github.com/sangtandoan/subscription_tracker/internal/utils"
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

func (h *subscriptionHandler) GetSubscriptionsBeforeNumDays(c *gin.Context) {
	numStr := c.Query("days")
	if numStr == "" {
		_ = c.Error(apperror.ErrInvalidJSON)
		return
	}

	num, err := strconv.Atoi(numStr)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res, err := h.s.GetSubscriptionsBeforeNumDays(c.Request.Context(), num)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("ok", res))
}

// GetAllSubscriptionsHandler godoc
//
//	@Summary		Get all subscriptions
//	@Description	get all subscriptions
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		service.GetAllSubscriptionsResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/subscriptions/ [get]
//
// This tells that this handler is protected by an API key
//
//	@Security		ApiKeyAuth
func (h *subscriptionHandler) GetAllSubscriptionsHandler(c *gin.Context) {
	req := &service.GetAllSubscriptionsRequest{}

	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}
	req.UserID = userID

	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		_ = c.Error(err)
		return
	}
	req.Limit = limitInt

	offset := c.Query("offset")
	if offset == "" {
		offset = "0"
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		_ = c.Error(err)
		return
	}
	req.Offset = offsetInt

	isCancelled := c.Query("is_cancelled")
	if isCancelled == "" {
		req.IsCancelled = nil
	} else {
		isCancelledBool, err := strconv.ParseBool(isCancelled)
		if err != nil {
			_ = c.Error(err)
			return
		}
		req.IsCancelled = &isCancelledBool
	}

	res, err := h.s.GetAllSubscriptions(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewAppResponse("getl all subscriptions successfully", res))
}

// CreateSubscriptionHandler godoc
//
//	@Summary		Create subscription
//	@Description	Create subscription
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		service.CreateSubscriptionRequest	true	"Create subscription request"
//	@Success		200				{array}		service.GetAllSubscriptionsResponse
//	@Failure		400				{object}	error
//	@Failure		404				{object}	error
//	@Failure		500				{object}	error
//	@Router			/subscriptions/ [post]
//
// This tells that this handler is protected by an API key
//
//	@Security		ApiKeyAuth
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

	userID, exists := c.Get(authenticator.SubClaim)
	if !exists {
		_ = c.Error(apperror.NewAppError(http.StatusUnauthorized, "can not find userID in context"))
		return
	}

	// uuid.UUID is not string
	// []byte is not string
	req.UserID, err = uuid.Parse(userID.(string))
	if err != nil {
		_ = c.Error(apperror.NewAppError(http.StatusUnauthorized, "can not parse to uuid"))
		return
	}

	res, err := h.s.CreateSubscription(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.NewAppResponse("created subscription sucessfully", res))
}
