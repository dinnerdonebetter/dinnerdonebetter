package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/lib/email"
	emailcfg "github.com/dinnerdonebetter/backend/internal/lib/email/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.EmailProberConfig]()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}
	cfg.Database.RunMigrations = false

	logger, tracerProvider, metricsProvider, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		log.Fatalf("could not establish observability pillars: %v", err)
	}

	ctx, span := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("email_prober_job")).StartSpan(ctx)
	defer span.End()

	emailer, err := emailcfg.ProvideEmailer(&cfg.Email, logger, tracerProvider, metricsProvider, tracing.BuildTracedHTTPClient())
	if err != nil {
		log.Fatalf("could not establish observability pillars: %v", err)
	}

	if err = emailer.SendEmail(ctx, &email.OutboundEmailMessage{
		ToAddress:   "verygoodsoftwarenotvirus@protonmail.com",
		ToName:      "Jeffrey",
		FromAddress: "email@dinnerdonebetter.dev",
		FromName:    "Testing",
		Subject:     "Testing",
		HTMLContent: "Hi",
	}); err != nil {
		log.Fatalf("could not send email: %v", err)
	}
}
