package main

import (
	"context"
	"fmt"
	"log"

	queuetest "github.com/dinnerdonebetter/backend/internal/build/jobs/queue_test"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func doTheThing(ctx context.Context) error {
	config.ConditionallyCease()

	cfg, err := config.LoadConfigFromEnvironment[config.QueueTestJobConfig]()
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	job, err := queuetest.Build(ctx, cfg)
	if err != nil {
		return fmt.Errorf("error building queue test job: %w", err)
	}

	if err = job.Do(ctx); err != nil {
		return fmt.Errorf("running queue test job: %w", err)
	}

	return nil
}

func main() {
	if err := doTheThing(context.Background()); err != nil {
		log.Fatal(err)
	}
}
