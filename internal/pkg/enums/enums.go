package enums

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
)

type Duration int

const (
	Monthly Duration = iota + 1
	SixMonths
	Yearly
)

var AllDurations = []string{"monthly", "6 months", "yearly"}

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
	fmt.Println(d.String())
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
