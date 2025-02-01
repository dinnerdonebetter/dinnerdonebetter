package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	searchdataindexscheduler "github.com/dinnerdonebetter/backend/internal/build/jobs/search_data_index_scheduler"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.SearchDataIndexSchedulerConfig]()
	if err != nil {
		log.Fatalf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	scheduler, err := searchdataindexscheduler.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("error building scheduler: %v", err)
	}

	if err = scheduler.IndexTypes(ctx); err != nil {
		log.Fatalf("error indexing types: %v", err)
	}
}
