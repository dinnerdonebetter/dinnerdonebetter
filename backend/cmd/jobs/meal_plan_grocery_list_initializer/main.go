package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	mealplangrocerylistinitializer "github.com/dinnerdonebetter/backend/internal/build/jobs/meal_plan_grocery_list_initializer"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func main() {
	ctx := context.Background()

	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.MealPlanGroceryListInitializerConfig]()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}
	cfg.Database.RunMigrations = false

	worker, err := mealplangrocerylistinitializer.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("error building mealplantaskcreator: %v", err)
	}

	if err = worker.Work(ctx); err != nil {
		log.Fatalf("error building mealplantaskcreator: %v", err)
	}
}
