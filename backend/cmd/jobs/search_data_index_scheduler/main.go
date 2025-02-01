package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	searchdataindexscheduler "github.com/dinnerdonebetter/backend/internal/build/jobs/search_data_index_scheduler"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func doTheThing(ctx context.Context) error {
	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.SearchDataIndexSchedulerConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	scheduler, err := searchdataindexscheduler.Build(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error building scheduler: %w", err)
	}

	if err = scheduler.IndexTypes(ctx); err != nil {
		return fmt.Errorf("error indexing types: %w", err)
	}

	return nil
}

func main() {
	if err := doTheThing(context.Background()); err != nil {
		log.Fatal(err)
	}
}
