package identifiers

import (
	"github.com/rs/xid"
)

// New produces a new string MealPlanTaskID.
func New() string {
	return xid.New().String()
}

// Validate validates a string MealPlanTaskID.
func Validate(x string) error {
	_, err := xid.FromString(x)
	return err
}
