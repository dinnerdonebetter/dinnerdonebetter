package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	tracing "github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	_ "go.uber.org/automaxprocs"
)

func doTheThing() error {
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	cfg, err := config.LoadConfigFromEnvironment[config.DBCleanerConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	logger, tracerProvider, metricsProvider, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		return fmt.Errorf("could not establish observability pillars: %w", err)
	}

	ctx, span := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("db_cleaner_job")).StartSpan(ctx)
	defer span.End()

	handledRecordsCounter, err := metricsProvider.NewInt64Counter("db_cleaner.handled_records")
	if err != nil {
		return fmt.Errorf("could not create metrics counter: %w", err)
	}

	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, tracerProvider, &cfg.Database)
	if err != nil {
		cancel()
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	cancel()
	defer dataManager.Close()

	deleted, err := dataManager.DeleteExpiredOAuth2ClientTokens(ctx)
	if err != nil {
		logger.Error("deleting expired oauth2 client tokens", err)
	}

	handledRecordsCounter.Add(ctx, deleted, metric.WithAttributes(
		attribute.KeyValue{
			Key:   "db_table",
			Value: attribute.StringValue("oauth2_clients"),
		},
	))

	return nil
}

func main() {
	log.Println("doing the thing")
	if err := doTheThing(); err != nil {
		log.Fatal(err)
	}
	log.Println("the thing is done")
}
