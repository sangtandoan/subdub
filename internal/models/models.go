package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/enums"
)

type User struct {
	CreatedAt time.Time
	ID        uuid.UUID
	Email     string
	Password  string
}

type Subscription struct {
	StartDate SubscriptionTime `json:"start_date"`
	EndDate   SubscriptionTime `json:"end_date"`
	Name      string           `json:"name,omitempty"`
	ID        uuid.UUID        `json:"id,omitempty"`
	UserID    uuid.UUID        `json:"user_id,omitempty"`
	Duration  enums.Duration   `json:"duration,omitempty"`
}

// create this type to enable marshal and unmarshal from format "YYYY-mm-dd"
type SubscriptionTime time.Time

func (st *SubscriptionTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	fmt.Println(s)

	// get rid of "" b/c the string will contains this
	s = strings.Trim(s, `"`)

	time, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*st = SubscriptionTime(time)
	return nil
}

func (st *SubscriptionTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(*st).Format("2006-01-02") + `"`), nil
}
