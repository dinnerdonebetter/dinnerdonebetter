package noop

import (
	"context"

	"github.com/prixfixeco/api_server/internal/email"
)

var _ email.Emailer = (*Emailer)(nil)

type (
	// Emailer doesn't send emails.
	Emailer struct{}
)

// NewNoopEmailer returns a new no-op Emailer.
func NewNoopEmailer() *Emailer {
	return &Emailer{}
}

// SendEmail sends an email.
func (e *Emailer) SendEmail(context.Context, *email.OutboundMessageDetails) error {
	return nil
}
