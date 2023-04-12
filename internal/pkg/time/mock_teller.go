package time

import (
	"time"

	"github.com/stretchr/testify/mock"
)

var _ Teller = (*StandardTimeTeller)(nil)

// MockTimeTeller is a mock time teller.
type MockTimeTeller struct {
	mock.Mock
}

// Now implements the Teller interface.
func (m *MockTimeTeller) Now() time.Time {
	return m.Called().Get(0).(time.Time)
}
