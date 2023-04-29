package email

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
)

var _ Emailer = (*NoopEmailer)(nil)

type (
	// NoopEmailer doesn't send emails.
	NoopEmailer struct {
		logger logging.Logger
	}
)

// NewNoopEmailer returns a new no-op NoopEmailer.
func NewNoopEmailer() (*NoopEmailer, error) {
	return &NoopEmailer{
		logger: zerolog.NewZerologLogger(logging.DebugLevel),
	}, nil
}

// SendEmail sends an email.
func (e *NoopEmailer) SendEmail(context.Context, *OutboundEmailMessage) error {
	e.logger.Info("NoopEmailer.SendEmail: no-op")
	return nil
}
