package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/models"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/enums"
)

type SubscriptionRepo interface {
	GetAllSubscriptions(ctx context.Context, userID uuid.UUID) ([]*SubscriptionRow, error)
	CreateSubscription(
		ctx context.Context,
		arg CreateSubscriptionParams,
	) (*SubscriptionRow, error)
	GetSubscriptionsBeforeNumDays(ctx context.Context, num int) ([]*SubscriptionRow, error)
}

type subscriptionRepo struct {
	db *sql.DB
}

func NewSubsciptionRepo(db *sql.DB) *subscriptionRepo {
	return &subscriptionRepo{db}
}

type SubscriptionRow struct {
	StartDate time.Time
	EndDate   time.Time
	Name      string
	Duration  string
	ID        uuid.UUID
	UserID    uuid.UUID
}

func (row *SubscriptionRow) MapToSubscriptionModel(sub *models.Subscription) error {
	var temp models.Subscription

	temp.ID = row.ID
	temp.UserID = row.UserID
	temp.Name = row.Name
	temp.StartDate = models.SubscriptionTime(row.StartDate)
	temp.EndDate = models.SubscriptionTime(row.EndDate)

	duration, err := enums.ParseString2Duration(row.Duration)
	if err != nil {
		return err
	}
	temp.Duration = duration

	*sub = temp
	return nil
}

func (repo *subscriptionRepo) GetAllSubscriptions(
	ctx context.Context,
	userID uuid.UUID,
) ([]*SubscriptionRow, error) {
	query := "SELECT id, user_id, name, start_date, end_date, duration FROM subscriptions WHERE user_id = $1"

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var res []*SubscriptionRow
	for rows.Next() {
		var sub SubscriptionRow
		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.Name,
			&sub.StartDate,
			&sub.EndDate,
			&sub.Duration,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, &sub)
	}

	return res, nil
}

type CreateSubscriptionParams struct {
	StartDate time.Time
	EndDate   time.Time
	Name      string
	ID        uuid.UUID
	UserID    uuid.UUID
	Duration  string
}

func (repo *subscriptionRepo) CreateSubscription(
	ctx context.Context,
	arg CreateSubscriptionParams,
) (*SubscriptionRow, error) {
	query := `
		INSERT INTO 
		subscriptions (id, user_id, name, start_date, end_date, duration) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, name, start_date, end_date, duration
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := repo.db.QueryRowContext(
		ctx,
		query,
		arg.ID,
		arg.UserID,
		arg.Name,
		time.Time(arg.StartDate),
		time.Time(arg.EndDate),
		arg.Duration,
	)

	var subcription SubscriptionRow
	err := row.Scan(
		&subcription.ID,
		&subcription.UserID,
		&subcription.Name,
		&subcription.StartDate,
		&subcription.EndDate,
		&subcription.Duration,
	)
	if err != nil {
		return nil, err
	}

	return &subcription, nil
}

func (repo *subscriptionRepo) GetSubscriptionsBeforeNumDays(
	ctx context.Context,
	num int,
) ([]*SubscriptionRow, error) {
	query := `
		SELECT id, user_id, name, start_date, end_date, duration 
		FROM subscriptions WHERE end_date <= $1 AND end_date + INTERVAL '1 day' >= $1
	`

	currentBeforeNumDays := time.Now().AddDate(0, 0, num)

	rows, err := repo.db.QueryContext(ctx, query, currentBeforeNumDays)
	if err != nil {
		return nil, err
	}

	var subs []*SubscriptionRow
	for rows.Next() {
		var sub SubscriptionRow
		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.Name,
			&sub.StartDate,
			&sub.EndDate,
			&sub.Duration,
		)
		if err != nil {
			return nil, err
		}

		subs = append(subs, &sub)
	}

	return subs, nil
}
