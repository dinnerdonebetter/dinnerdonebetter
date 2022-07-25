package email

import (
	"context"
	"github.com/prixfixeco/api_server/pkg/types"

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
func (m *MockEmailer) SendEmail(ctx context.Context, details *OutboundEmailMessage) error {
	return m.Called(ctx, details).Error(0)
}

func (m *MockEmailer) SendHouseholdInvitationEmail(ctx context.Context, householdInvitation *types.HouseholdInvitation) error {
	return m.Called(ctx, householdInvitation).Error(0)
}
