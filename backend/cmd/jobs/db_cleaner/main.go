package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	dbcleaner "github.com/dinnerdonebetter/backend/internal/build/jobs/db_cleaner"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"

	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.DBCleanerConfig]()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}
	cfg.Database.RunMigrations = false

	logger, err := cfg.Observability.Logging.ProvideLogger(ctx)
	if err != nil {
		log.Fatalf("error getting logger: %v", err)
	}

	dbCleaner, err := dbcleaner.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("error building db cleaner: %v", err)
	}

	if err = dbCleaner.Do(ctx); err != nil {
		observability.AcknowledgeError(err, logger, nil, "cleaning database")
	}
}
