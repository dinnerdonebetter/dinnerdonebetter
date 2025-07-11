package main

import (
	"context"
	"fmt"
	"log"

	mealplantaskcreator "github.com/dinnerdonebetter/backend/internal/build/jobs/meal_plan_task_creator"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func doTheThing(ctx context.Context) error {
	config.ConditionallyCease()

	cfg, err := config.LoadConfigFromEnvironment[config.MealPlanTaskCreatorConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}
	cfg.Database.RunMigrations = false

	worker, err := mealplantaskcreator.Build(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error building mealplantaskcreator: %w", err)
	}

	if err = worker.Work(ctx); err != nil {
		return fmt.Errorf("error building mealplantaskcreator: %w", err)
	}

	return nil
}

func main() {
	if err := doTheThing(context.Background()); err != nil {
		log.Fatal(err)
	}
}
