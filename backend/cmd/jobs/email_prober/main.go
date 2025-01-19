package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/email"
	emailcfg "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"go.opentelemetry.io/otel"
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

	logger, err := cfg.Observability.Logging.ProvideLogger(ctx)
	if err != nil {
		return fmt.Errorf("could not create logger: %w", err)
	}

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error("initializing tracer", initializeTracerErr)
	}
	otel.SetTracerProvider(tracerProvider)

	metricsProvider, err := cfg.Observability.Metrics.ProvideMetricsProvider(ctx, logger)
	if err != nil {
		logger.Error("initializing metrics provider", err)
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
