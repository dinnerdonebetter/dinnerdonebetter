package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	dbcleaner "github.com/dinnerdonebetter/backend/internal/build/jobs/db_cleaner"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func doTheThing(ctx context.Context) error {
	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.DBCleanerConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	dbCleaner, err := dbcleaner.Build(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error building db cleaner: %w", err)
	}

	if err = dbCleaner.Do(ctx); err != nil {
		return fmt.Errorf("cleaning database: %w", err)
	}

	return nil
}

func main() {
	if err := doTheThing(context.Background()); err != nil {
		log.Fatal(err)
	}
}
