package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/build/jobs/meal_plan_task_creator"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.MealPlanTaskCreatorConfig]()
	if err != nil {
		log.Fatalf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	worker, err := mealplantaskcreator.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("error building mealplantaskcreator: %w", err)
	}

	if err = worker.Work(ctx); err != nil {
		log.Fatalf("error building mealplantaskcreator: %w", err)
	}
}
