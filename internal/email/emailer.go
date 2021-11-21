package email

import (
	"context"
)

type (
	// APIToken is used to authenticate an email service.
	APIToken string

	// OutboundMessageDetails is a collection of fields that are useful for sending emails.
	OutboundMessageDetails struct {
		ToAddress   string
		ToName      string
		FromAddress string
		FromName    string
		Subject     string
		HTMLContent string
	}

	// Emailer represents a service that can send emails.
	Emailer interface {
		SendEmail(ctx context.Context, details *OutboundMessageDetails) error
	}
)
