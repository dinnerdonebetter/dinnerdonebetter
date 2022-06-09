package email

import (
	"fmt"
	"os"

	"github.com/prixfixeco/api_server/pkg/types"
)

var urlMap = map[string]string{
	"dev": "https://www.prixfixe.dev",
}

// BuildInviteMemberEmail builds an email notifying a user that they've been invited to join a household.
func BuildInviteMemberEmail(householdInvitation *types.HouseholdInvitation) (*OutboundMessageDetails, error) {
	env := os.Getenv("PF_ENVIRONMENT")
	envAddr, ok := urlMap[env]
	if !ok {
		return nil, fmt.Errorf("no available URL for")
	}

	msg := &OutboundMessageDetails{
		ToAddress:   householdInvitation.ToEmail,
		ToName:      "",
		FromAddress: "invites@prixfixe.email",
		FromName:    "PrixFixe",
		Subject:     "You've been invited to join a household on PrixFixe!",
		HTMLContent: fmt.Sprintf(`Register at: %s/register?token=%s&destination=%s`, envAddr, householdInvitation.Token, householdInvitation.DestinationHousehold),
	}

	return msg, nil
}
