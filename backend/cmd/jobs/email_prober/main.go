package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/lib/email"
	emailcfg "github.com/dinnerdonebetter/backend/internal/lib/email/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	_ "go.uber.org/automaxprocs"
)

// I'm not going to bother with turning this into a proper worker because it probably just should not exist.
func doTheThing(ctx context.Context) error {
	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.EmailProberConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	logger, tracerProvider, metricsProvider, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		return fmt.Errorf("could not establish observability pillars: %w", err)
	}

	ctx, span := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("email_prober_job")).StartSpan(ctx)
	defer span.End()

	//nolint:contextcheck // the context it's complaining about is irrelevant.
	emailer, err := emailcfg.ProvideEmailer(&cfg.Email, logger, tracerProvider, metricsProvider, tracing.BuildTracedHTTPClient())
	if err != nil {
		return fmt.Errorf("could not establish observability pillars: %w", err)
	}

	if err = emailer.SendEmail(ctx, &email.OutboundEmailMessage{
		ToAddress:   "verygoodsoftwarenotvirus@protonmail.com",
		ToName:      "Jeffrey",
		FromAddress: "email@dinnerdonebetter.dev",
		FromName:    "Testing",
		Subject:     "Testing",
		HTMLContent: "Hi",
	}); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}

	return nil
}

func main() {
	if err := doTheThing(context.Background()); err != nil {
		log.Fatal(err)
	}
}
