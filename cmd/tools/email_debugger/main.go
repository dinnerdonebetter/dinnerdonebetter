package main

import (
	"context"
	"net/http"
	"os"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/prixfixeco/api_server/internal/email/sendgrid"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	apiToken         string
	destinationEmail string
	webappURL        string
)

func init() {
	flag.StringVarP(&apiToken, "api-token", "a", "", "SendGrid API token")
	flag.StringVarP(&destinationEmail, "to-email", "t", "", "email address to test sending emails to")
	flag.StringVarP(&webappURL, "webapp-url", "u", "https://www.prixfixe.dev", "webapp URL to point users towards")
}

func main() {
	ctx := context.Background()
	logger := zerolog.NewZerologLogger()

	if err := os.Setenv("PF_ENVIRONMENT", "dev"); err != nil {
		panic(err)
	}

	cfg := sendgrid.Config{
		APIToken:                            apiToken,
		WebAppURL:                           "https://www.prixfixe.fake.lol",
		HouseholdInviteOutboundEmailAddress: destinationEmail,
	}

	emailer, err := sendgrid.NewSendGridEmailer(
		cfg,
		logger,
		tracing.NewNoopTracerProvider(),
		&http.Client{Timeout: 5 * time.Second},
	)
	if err != nil {
		panic(err)
	}

	householdInvitation := &types.HouseholdInvitation{
		ToEmail:              destinationEmail,
		Token:                "blah_example_token_blah",
		DestinationHousehold: types.Household{ID: "__te$ting__"},
	}

	if sendErr := emailer.SendHouseholdInvitationEmail(ctx, householdInvitation); sendErr != nil {
		panic(sendErr)
	}

	println("yay")
}
