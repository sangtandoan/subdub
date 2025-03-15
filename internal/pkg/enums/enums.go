package enums

import (
	"net/http"
	"strings"
	"time"

	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
)

type Duration int

const (
	Weekly Duration = iota + 1
	Monthly
	SixMonths
	Yearly
)

var AllDurations = []string{"weekly", "monthly", "6 months", "yearly"}

func (d *Duration) String() string {
	return AllDurations[*d-1]
}

func ParseString2Duration(s string) (Duration, error) {
	for i, duration := range AllDurations {
		if s == duration {
			return Duration(i + 1), nil
		}
	}

	return Duration(0), apperror.NewAppError(http.StatusBadRequest, "invalid duration")
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// this is also validate for duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, `"`)

	duration, err := ParseString2Duration(s)
	if err != nil {
		return err
	}

	*d = duration

	return nil
}

func (d *Duration) AddDurationToTime(start time.Time) time.Time {
	var end time.Time

	switch *d {
	case Weekly:
		end = time.Time(start).AddDate(0, 0, 7)
	case Monthly:
		end = time.Time(start).AddDate(0, 1, 0)
	case SixMonths:
		end = time.Time(start).AddDate(0, 6, 0)
	case Yearly:
		end = time.Time(start).AddDate(1, 0, 0)
	}

	return end
}
