package email

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Emailer = (*MockEmailer)(nil)

type (
	// MockEmailer represents a service that can send emails.
	MockEmailer struct {
		mock.Mock
	}
)

// SendEmail is a mock function.
func (m *MockEmailer) SendEmail(ctx context.Context, details *OutboundMessageDetails) error {
	return m.Called(ctx, details).Error(0)
}
