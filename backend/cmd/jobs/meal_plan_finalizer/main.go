package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	mealplanfinalizer "github.com/dinnerdonebetter/backend/internal/build/jobs/meal_plan_finalizer"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func doTheThing(ctx context.Context) error {
	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.MealPlanFinalizerConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	worker, err := mealplanfinalizer.Build(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error building mealplantaskcreator: %w", err)
	}

	if _, err = worker.Work(ctx); err != nil {
		return fmt.Errorf("error building mealplantaskcreator: %w", err)
	}

	return nil
}

func main() {
	if err := doTheThing(context.Background()); err != nil {
		log.Fatal(err)
	}
}
