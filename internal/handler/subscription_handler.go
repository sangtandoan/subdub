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

func (h *subscriptionHandler) GetAllSubscriptionsHandler(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.Error(err)
		return
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		_ = c.Error(err)
		return
	}

	offset := c.Query("offset")
	if offset == "" {
		offset = "0"
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		_ = c.Error(err)
		return
	}

	req := service.GetAllSubscriptionsRequest{
		UserID: userID,
		Limit:  limitInt,
		Offset: offsetInt,
	}

	res, err := h.s.GetAllSubscriptions(c.Request.Context(), &req)
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
