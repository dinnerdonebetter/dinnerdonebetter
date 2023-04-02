package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/prixfixeco/backend/internal/database"
	dbconfig "github.com/prixfixeco/backend/internal/database/config"
	"github.com/prixfixeco/backend/internal/database/postgres"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/logging/zap"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/pointers"
	"github.com/prixfixeco/backend/internal/search"
	"github.com/prixfixeco/backend/internal/search/algolia"
	searchcfg "github.com/prixfixeco/backend/internal/search/config"
	"github.com/prixfixeco/backend/pkg/types"
)

func main() {
	ctx := context.Background()
	logger := zap.NewZapLogger(logging.DebugLevel)
	tracerProvider := tracing.NewNoopTracerProvider()

	cfg := &searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			AppID:  os.Getenv("ALGOLIA_APP_ID"),
			APIKey: os.Getenv("ALGOLIA_API_KEY"),
		},
	}

	im, err := searchcfg.ProvideIndexManager[types.ValidPreparation](ctx, logger, tracerProvider, cfg, "valid_preparations")
	if err != nil {
		log.Fatal(fmt.Errorf("initializing index manager: %w", err))
	}

	if err = wipeIndex(ctx, im); err != nil {
		log.Fatal(fmt.Errorf("wiping index: %w", err))
	}

	dbConfig := &dbconfig.Config{
		ConnectionDetails: database.ConnectionDetails(os.Getenv("DATABASE_URL")),
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, dbConfig, tracing.NewNoopTracerProvider())
	if err != nil {
		log.Fatal(fmt.Errorf("initializing database client: %w", err))
	}

	filter := types.DefaultQueryFilter()
	filter.Limit = pointers.Uint8(50)

	var thresholdMet bool
	for !thresholdMet {
		var preparations *types.QueryFilteredResult[types.ValidPreparation]
		preparations, err = dataManager.GetValidPreparations(ctx, filter)
		if err != nil {
			log.Fatal(fmt.Errorf("getting valid preparations: %w", err))
		}

		for _, prep := range preparations.Data {
			if err = im.Index(ctx, prep.ID, prep); err != nil {
				log.Fatal(fmt.Errorf("indexing preparation: %w", err))
			}
		}

		thresholdMet = len(preparations.Data) == 0
		*filter.Page++
	}

	results, err := im.Search(ctx, "m")
	if err != nil {
		log.Fatal(fmt.Errorf("searching index: %w", err))
	}

	log.Println(results)
}

func wipeIndex(ctx context.Context, im search.IndexManager[types.ValidPreparation]) error {
	if err := im.Wipe(ctx); err != nil {
		return fmt.Errorf("deleting index: %w", err)
	}

	return nil
}
