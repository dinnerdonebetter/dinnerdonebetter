package time

import "time"

// Teller is a time teller.
type Teller interface {
	Now() time.Time
}

var _ Teller = (*StandardTimeTeller)(nil)

// StandardTimeTeller is the standard library time teller.
type StandardTimeTeller struct{}

// Now implements the Teller interface.
func (t *StandardTimeTeller) Now() time.Time {
	return time.Now()
}
