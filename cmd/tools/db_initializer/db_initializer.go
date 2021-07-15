package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"gitlab.com/prixfixe/prixfixe/internal/config"
	"gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"

	flag "github.com/spf13/pflag"
)

var (
	connectionString string
	debug            bool
)

func init() {
	flag.StringVarP(&connectionString, "url", "u", "", "where the target instance is hosted")
	flag.BoolVarP(&debug, "debug", "d", false, "whether debug mode is enabled")
}

func main() {
	flag.Parse()

	ctx := context.Background()
	logger := logging.ProvideLogger(logging.Config{Provider: logging.ProviderZerolog})

	if connectionString == "" {
		logger.Fatal(errors.New("connection string must not be empty"))
	}

	cfg := &config.InstanceConfig{
		Database: dbconfig.Config{
			Provider:          dbconfig.PostgresProvider,
			ConnectionDetails: database.ConnectionDetails(connectionString),
		},
	}

	db, err := dbconfig.ProvideDatabaseConnection(logger, &cfg.Database)
	if err != nil {
		logger.Fatal(fmt.Errorf("initializing client: %w", err))
	}

	client, err := config.ProvideDatabaseClient(ctx, logger, db, cfg)
	if err != nil {
		logger.Fatal(fmt.Errorf("initializing client: %w", err))
	}

	logger.Debug("initialized db client")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for _, instrument := range validInstruments {
			if _, instrumentCreationErr := client.CreateValidInstrument(ctx, instrument, 1); instrumentCreationErr != nil {
				logger.Error(instrumentCreationErr, "creating valid instrument")
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for _, preparation := range validPreparations {
			if _, preparationCreationErr := client.CreateValidPreparation(ctx, preparation, 1); preparationCreationErr != nil {
				logger.Error(preparationCreationErr, "creating valid instrument")
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for _, ingredient := range validIngredients {
			if _, instrumentCreationErr := client.CreateValidIngredient(ctx, ingredient, 1); instrumentCreationErr != nil {
				logger.Error(instrumentCreationErr, "creating valid ingredient")
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
