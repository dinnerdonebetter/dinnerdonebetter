/*
Command index_initializer is a CLI that takes in some data via flags about your
database and the type you want to index, and hydrates a Bleve index full of that type.
This tool is to be used in the event of some data corruption that takes the search index
out of commission.
*/
package main

import (
	"context"
	"log"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	"gitlab.com/prixfixe/prixfixe/internal/v1/search"
	"gitlab.com/prixfixe/prixfixe/internal/v1/search/bleve"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	flag "github.com/spf13/pflag"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/zerolog"
)

var (
	indexOutputPath string
	typeName        string

	dbConnectionDetails string
	databaseType        string

	deadline time.Duration

	validTypeNames = map[string]struct{}{
		"valid_instruments":  {},
		"valid_ingredients":  {},
		"valid_preparations": {},
	}

	validDatabaseTypes = map[string]struct{}{
		config.PostgresProviderKey: {},
	}
)

const (
	outputPathVerboseFlagName   = "output"
	dbConnectionVerboseFlagName = "db_connection"
	dbTypeVerboseFlagName       = "db_type"
)

func init() {
	flag.StringVarP(&indexOutputPath, outputPathVerboseFlagName, "o", "", "output path for bleve index")
	flag.StringVarP(&typeName, "type", "t", "", "which type to create bleve index for")

	flag.StringVarP(&dbConnectionDetails, dbConnectionVerboseFlagName, "c", "", "connection string for the relevant database")
	flag.StringVarP(&databaseType, dbTypeVerboseFlagName, "b", "", "which type of database to connect to")

	flag.DurationVarP(&deadline, "deadline", "d", time.Minute, "amount of time to spend adding to the index")
}

func main() {
	flag.Parse()
	logger := zerolog.NewZeroLogger().WithName("search_index_initializer")
	ctx := context.Background()

	if indexOutputPath == "" {
		log.Fatalf("No output path specified, please provide one via the --%s flag", outputPathVerboseFlagName)
	} else if dbConnectionDetails == "" {
		log.Fatalf("No database connection details %q specified, please provide one via the --%s flag", dbConnectionDetails, dbConnectionVerboseFlagName)
	} else if _, ok := validTypeNames[typeName]; !ok {
		log.Fatalf("Invalid type name %q specified, one of [ 'valid_instrument', 'valid_ingredient', 'valid_preparation' ] expected", typeName)
	} else if _, ok := validDatabaseTypes[databaseType]; !ok {
		log.Fatalf("Invalid database type %q specified, please provide one via the --%s flag", databaseType, dbTypeVerboseFlagName)
	}

	im, err := bleve.NewBleveIndexManager(search.IndexPath(indexOutputPath), search.IndexName(typeName), logger)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &config.ServerConfig{
		Database: config.DatabaseSettings{
			Provider:          databaseType,
			ConnectionDetails: database.ConnectionDetails(dbConnectionDetails),
		},
		Metrics: config.MetricsSettings{
			DBMetricsCollectionInterval: time.Second,
		},
	}

	// connect to our database.
	logger.Debug("connecting to database")
	rawDB, err := cfg.ProvideDatabaseConnection(logger)
	if err != nil {
		log.Fatalf("error establishing connection to database: %v", err)
	}

	// establish the database client.
	logger.Debug("setting up database client")
	dbClient, err := cfg.ProvideDatabaseClient(ctx, logger, rawDB, false)
	if err != nil {
		log.Fatalf("error initializing database client: %v", err)
	}

	switch typeName {
	case "valid_instruments":
		outputChan := make(chan []models.ValidInstrument)
		if queryErr := dbClient.GetAllValidInstruments(ctx, outputChan); queryErr != nil {
			log.Fatalf("error fetching valid instruments from database: %v", err)
		}

		for {
			select {
			case validInstruments := <-outputChan:
				for _, x := range validInstruments {
					if searchIndexErr := im.Index(ctx, x.ID, x); searchIndexErr != nil {
						logger.WithValue("id", x.ID).Error(searchIndexErr, "error adding to search index")
					}
				}
			case <-time.After(deadline):
				logger.Info("terminating")
				return
			}
		}
	case "valid_ingredients":
		outputChan := make(chan []models.ValidIngredient)
		if queryErr := dbClient.GetAllValidIngredients(ctx, outputChan); queryErr != nil {
			log.Fatalf("error fetching valid ingredients from database: %v", err)
		}

		for {
			select {
			case validIngredients := <-outputChan:
				for _, x := range validIngredients {
					if searchIndexErr := im.Index(ctx, x.ID, x); searchIndexErr != nil {
						logger.WithValue("id", x.ID).Error(searchIndexErr, "error adding to search index")
					}
				}
			case <-time.After(deadline):
				logger.Info("terminating")
				return
			}
		}
	case "valid_preparations":
		outputChan := make(chan []models.ValidPreparation)
		if queryErr := dbClient.GetAllValidPreparations(ctx, outputChan); queryErr != nil {
			log.Fatalf("error fetching valid preparations from database: %v", err)
		}

		for {
			select {
			case validPreparations := <-outputChan:
				for _, x := range validPreparations {
					if searchIndexErr := im.Index(ctx, x.ID, x); searchIndexErr != nil {
						logger.WithValue("id", x.ID).Error(searchIndexErr, "error adding to search index")
					}
				}
			case <-time.After(deadline):
				logger.Info("terminating")
				return
			}
		}
	default:
		log.Fatal("this should never occur")
	}
}
