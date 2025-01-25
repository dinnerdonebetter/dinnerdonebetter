package email

import (
	"context"
)

type (
	// APIToken is used to authenticate an email service.
	APIToken string

	// OutboundEmailMessage is a collection of fields that are useful for sending emails.
	OutboundEmailMessage struct {
		ToAddress   string
		ToName      string
		FromAddress string
		FromName    string
		Subject     string
		HTMLContent string
	}

	// Emailer represents a service that can send emails.
	Emailer interface {
		SendEmail(ctx context.Context, details *OutboundEmailMessage) error
	}
)
