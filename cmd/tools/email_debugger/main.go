package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/email/sendgrid"
	"github.com/prixfixeco/backend/internal/observability"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zerolog"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types"

	flag "github.com/spf13/pflag"
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
	logger := zerolog.NewZerologLogger(logging.DebugLevel)

	if err := os.Setenv("PF_ENVIRONMENT", "dev"); err != nil {
		panic(err)
	}

	cfg := &sendgrid.Config{
		APIToken:  apiToken,
		WebAppURL: "https://www.prixfixe.fake.lol",
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

	msg, emailGenerationErr := email.BuildInviteMemberEmail(householdInvitation)
	if emailGenerationErr != nil {
		panic(observability.PrepareError(emailGenerationErr, nil, "building email message"))
	}

	if err = emailer.SendEmail(ctx, msg); err != nil {
		panic(observability.PrepareError(err, nil, "sending email notice"))
	}

	println("yay")
}
