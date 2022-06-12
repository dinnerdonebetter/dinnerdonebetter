package email

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"

	"github.com/prixfixeco/api_server/pkg/types"
)

var urlMap = map[string]string{
	"dev":     "https://www.prixfixe.dev",
	"testing": "https://not.real.lol",
}

var (
	//go:embed templates/invite.tmpl
	outgoingInviteTemplate string
)

type inviteContent struct {
	WebAppURL            string
	Token                string
	DestinationHousehold string
}

// BuildInviteMemberEmail builds an email notifying a user that they've been invited to join a household.
func BuildInviteMemberEmail(householdInvitation *types.HouseholdInvitation) (*OutboundMessageDetails, error) {
	env := os.Getenv("PF_ENVIRONMENT")
	if env == "" {
		env = "testing"
	}

	envAddr, ok := urlMap[env]
	if !ok {
		return nil, fmt.Errorf("no available URL for")
	}

	content := &inviteContent{
		WebAppURL:            envAddr,
		Token:                householdInvitation.Token,
		DestinationHousehold: householdInvitation.DestinationHousehold,
	}

	tmpl := template.Must(template.New("").Funcs(map[string]any{}).Parse(outgoingInviteTemplate))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, content); err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &OutboundMessageDetails{
		ToAddress:   householdInvitation.ToEmail,
		ToName:      "",
		FromAddress: "invites@prixfixe.dev",
		FromName:    "PrixFixe",
		Subject:     "You've been invited to join a household on PrixFixe!",
		HTMLContent: b.String(),
	}

	return msg, nil
}
