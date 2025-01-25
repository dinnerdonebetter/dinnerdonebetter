package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/lib/email"
	emailcfg "github.com/dinnerdonebetter/backend/internal/lib/email/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	tracing "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	_ "go.uber.org/automaxprocs"
)

func doTheThing() error {
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		return nil
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

	emailer, err := emailcfg.ProvideEmailer(&cfg.Email, logger, tracerProvider, metricsProvider, tracing.BuildTracedHTTPClient())
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring outbound emailer")
	}

	return emailer.SendEmail(ctx, &email.OutboundEmailMessage{
		ToAddress:   "verygoodsoftwarenotvirus@protonmail.com",
		ToName:      "Jeffrey",
		FromAddress: "email@dinnerdonebetter.dev",
		FromName:    "Testing",
		Subject:     "Testing",
		HTMLContent: "Hi",
	})
}

func main() {
	log.Println("doing the thing")
	if err := doTheThing(); err != nil {
		log.Fatal(err)
	}
	log.Println("the thing is done")
}
