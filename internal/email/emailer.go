package email

import (
	"context"
	"github.com/prixfixeco/api_server/pkg/types"
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
		SendHouseholdInvitationEmail(ctx context.Context, householdInvitation *types.HouseholdInvitation) error
	}
)
