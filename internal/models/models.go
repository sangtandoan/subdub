package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/enums"
)

type User struct {
	CreatedAt time.Time
	Email     string
	Password  string
	ID        uuid.UUID
}

type Subscription struct {
	StartDate SubscriptionTime `json:"start_date"`
	EndDate   SubscriptionTime `json:"end_date"`
	Name      string           `json:"name,omitempty"`
	ID        uuid.UUID        `json:"id,omitempty"`
	UserID    uuid.UUID        `json:"user_id,omitempty"`
	Duration  enums.Duration   `json:"duration,omitempty"`
}

type Session struct {
	CreatedAt    time.Time
	ExpiresAt    time.Time
	RefreshToken string
	UserEmail    string
	ID           uuid.UUID
	IsRevoked    bool
}

type AuthProvider struct {
	CreatedAt  time.Time
	Provider   string
	ProviderID string
	ID         uuid.UUID
	UserID     uuid.UUID
}

// create this type to enable marshal and unmarshal from format "YYYY-mm-dd"
// if using normal time.Time, when unmarshal will occur error
type SubscriptionTime time.Time

func (st *SubscriptionTime) UnmarshalJSON(b []byte) error {
	s := string(b)

	// get rid of "" b/c the string will contains this
	s = strings.Trim(s, `"`)

	// string -> time : time.Parse
	// time -> string: time.Format
	normalTime, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*st = SubscriptionTime(normalTime)
	return nil
}

func (st *SubscriptionTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*st).Format("2006-01-02") + `"`), nil
}
