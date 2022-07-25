package email

import (
	"context"
)

var _ Emailer = (*NoopEmailer)(nil)

type (
	// NoopEmailer doesn't send emails.
	NoopEmailer struct{}
)

// NewNoopEmailer returns a new no-op NoopEmailer.
func NewNoopEmailer() (*NoopEmailer, error) {
	return &NoopEmailer{}, nil
}

// SendEmail sends an email.
func (e *NoopEmailer) SendEmail(context.Context, *OutboundEmailMessage) error {
	return nil
}
