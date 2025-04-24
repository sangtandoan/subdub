package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/models"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/enums"
)

type SubscriptionRepo interface {
	GetAllSubscriptions(
		ctx context.Context,
		arg *GetAllSubscriptionsParams,
	) ([]*SubscriptionRow, int, error)
	CreateSubscription(
		ctx context.Context,
		arg CreateSubscriptionParams,
	) (*SubscriptionRow, error)
	GetSubscriptionsBeforeNumDays(ctx context.Context, num int) ([]*SubscriptionRow, error)
	GetSubscriptionsNeedUpdateStartAndEndDate(ctx context.Context) ([]*SubscriptionRow, error)
	UpdateSubscriptionStartAndEndDate(
		ctx context.Context,
		arg *UpdateSubscriptionStartAndEndDateParams,
	) error
}

type subscriptionRepo struct {
	db *sql.DB
}

func NewSubsciptionRepo(db *sql.DB) *subscriptionRepo {
	return &subscriptionRepo{db}
}

// We use this struct to scan the result from the database
// because models.Subscription has a custom type SubscriptionTime
// which postgres driver can not scan directly to it
type SubscriptionRow struct {
	StartDate   time.Time
	EndDate     time.Time
	Name        string
	Duration    string
	ID          uuid.UUID
	UserID      uuid.UUID
	IsCancelled bool
}

func (row *SubscriptionRow) MapToSubscriptionModel(sub *models.Subscription) error {
	var temp models.Subscription

	temp.ID = row.ID
	temp.UserID = row.UserID
	temp.Name = row.Name
	temp.IsCancelled = row.IsCancelled
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

type GetAllSubscriptionsParams struct {
	IsCancelled *bool
	UserID      uuid.UUID
	Limit       int
	Offset      int
}

func (repo *subscriptionRepo) GetAllSubscriptions(
	ctx context.Context,
	arg *GetAllSubscriptionsParams,
) ([]*SubscriptionRow, int, error) {
	query := `
		SELECT id, user_id, name, start_date, end_date, duration, is_cancelled FROM subscriptions`

	var whereClauses []string
	var args []any
	argIndex := 1

	// optinal query param is_cancelled
	if arg.IsCancelled != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_cancelled = $%d", argIndex))
		args = append(args, *arg.IsCancelled)
		argIndex++
	}

	// add all optionals to query string
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// add pagination to quer string
	query += fmt.Sprintf(" ORDER BY start_date ASC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, arg.Limit, arg.Offset)
	fmt.Println(query)

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
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
			&sub.IsCancelled,
		)
		if err != nil {
			return nil, 0, err
		}

		res = append(res, &sub)
	}

	query = `SELECT COUNT(*) FROM subscriptions WHERE user_id = $1`
	// This will create a new memory allocation for args
	// => costs more memory
	// args = []any{}

	// This will clear the slice but keep the memory allocation
	// => costs less memory
	args = args[:0]
	args = append(args, arg.UserID)

	if arg.IsCancelled != nil {
		query += " AND is_cancelled = $2"
		args = append(args, *arg.IsCancelled)
	}

	var count int
	err = repo.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return res, count, nil
}

type CreateSubscriptionParams struct {
	StartDate time.Time
	EndDate   time.Time
	Name      string
	Duration  string
	ID        uuid.UUID
	UserID    uuid.UUID
}

func (repo *subscriptionRepo) CreateSubscription(
	ctx context.Context,
	arg CreateSubscriptionParams,
) (*SubscriptionRow, error) {
	query := `
		INSERT INTO 
		subscriptions (id, user_id, name, start_date, end_date, duration) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, name, start_date, end_date, duration, is_cancelled
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := repo.db.QueryRowContext(
		ctx,
		query,
		arg.ID,
		arg.UserID,
		arg.Name,
		arg.StartDate,
		arg.EndDate,
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
		&subcription.IsCancelled,
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
		SELECT id, user_id, name, start_date, end_date, duration, is_cancelled
		FROM subscriptions WHERE end_date <= $1 AND end_date + INTERVAL '1 day' >= $1
	`

	futureAfterNumDays := time.Now().AddDate(0, 0, num)

	rows, err := repo.db.QueryContext(ctx, query, futureAfterNumDays)
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
			&sub.IsCancelled,
		)
		if err != nil {
			return nil, err
		}

		subs = append(subs, &sub)
	}

	return subs, nil
}

func (repo *subscriptionRepo) GetSubscriptionsNeedUpdateStartAndEndDate(
	ctx context.Context,
) ([]*SubscriptionRow, error) {
	query := `
	    SELECT id, user_id, name, start_date, end_date, duration, is_cancelled 
		FROM subscriptions
	    WHERE end_date <= $1 and end_date + INTERVAL '1 day' >= $1
	`
	now := time.Now()
	rows, err := repo.db.QueryContext(ctx, query, now)
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
			&sub.IsCancelled,
		)
		if err != nil {
			return nil, err
		}

		subs = append(subs, &sub)
	}

	return subs, nil
}

type UpdateSubscriptionStartAndEndDateParams struct {
	StartDate time.Time
	EndDate   time.Time
	ID        uuid.UUID
}

func (reop *subscriptionRepo) UpdateSubscriptionStartAndEndDate(
	ctx context.Context,
	arg *UpdateSubscriptionStartAndEndDateParams,
) error {
	query := `
		UPDATE subscriptions 
		SET start_date = $1, end_date = $2
		WHERE id = $3
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := reop.db.ExecContext(ctx, query, arg.StartDate, arg.EndDate, arg.ID)
	if err != nil {
		return err
	}

	return nil
}
