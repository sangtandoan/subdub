package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/models"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/enums"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
)

type SubscriptionService interface {
	GetAllSubscriptions(ctx context.Context, userID uuid.UUID) ([]*models.Subscription, error)
	CreateSubscription(
		ctx context.Context,
		req *CreateSubscriptionRequest,
	) (*models.Subscription, error)
	GetSubscriptionsBeforeNumDays(
		ctx context.Context,
		num int,
	) ([]*models.Subscription, error)
}

type subscriptionService struct {
	repo repo.SubscriptionRepo
}

func NewSubscriptionService(repo repo.SubscriptionRepo) *subscriptionService {
	return &subscriptionService{repo}
}

func (s *subscriptionService) GetSubscriptionsBeforeNumDays(
	ctx context.Context,
	num int,
) ([]*models.Subscription, error) {
	subs, err := s.repo.GetSubscriptionsBeforeNumDays(ctx, num)
	if err != nil {
		return nil, err
	}

	var arr []*models.Subscription
	for _, row := range subs {
		var sub models.Subscription
		err := row.MapToSubscriptionModel(&sub)
		if err != nil {
			return nil, err
		}

		arr = append(arr, &sub)
	}

	return arr, nil
}

func (s *subscriptionService) GetAllSubscriptions(
	ctx context.Context,
	userID uuid.UUID,
) ([]*models.Subscription, error) {
	res, err := s.repo.GetAllSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var arr []*models.Subscription
	for _, row := range res {
		var sub models.Subscription
		err := row.MapToSubscriptionModel(&sub)
		if err != nil {
			return nil, err
		}

		arr = append(arr, &sub)
	}

	return arr, nil
}

type CreateSubscriptionRequest struct {
	StartDate models.SubscriptionTime `json:"start_date"         validate:"required"`
	Name      string                  `json:"name,omitempty"     validate:"required,min=3,max=20"`
	UserID    uuid.UUID               `json:"-"                  validate:"-"`
	Duration  enums.Duration          `json:"duration,omitempty" validate:"required"`
}

func (s *subscriptionService) CreateSubscription(
	ctx context.Context,
	req *CreateSubscriptionRequest,
) (*models.Subscription, error) {
	endDate := calculateEndDate(req.StartDate, req.Duration)

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	arg := repo.CreateSubscriptionParams{
		ID:        id,
		UserID:    req.UserID,
		StartDate: time.Time(req.StartDate),
		EndDate:   endDate,
		Name:      req.Name,
		Duration:  req.Duration.String(),
	}
	row, err := s.repo.CreateSubscription(ctx, arg)
	if err != nil {
		return nil, err
	}

	var res models.Subscription
	err = row.MapToSubscriptionModel(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func calculateEndDate(
	startDate models.SubscriptionTime,
	duration enums.Duration,
) time.Time {
	return duration.AddDurationToTime(time.Time(startDate))
}
