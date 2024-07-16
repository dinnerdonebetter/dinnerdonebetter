package emailmock

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/email"

	"github.com/stretchr/testify/mock"
)

var _ email.Emailer = (*Emailer)(nil)

type (
	// Emailer represents a service that can send emails.
	Emailer struct {
		mock.Mock
	}
)

// SendEmail is a mock function.
func (m *Emailer) SendEmail(ctx context.Context, details *email.OutboundEmailMessage) error {
	return m.Called(ctx, details).Error(0)
}
