package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dinnerdonebetter/backend/internal/database"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zap"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/internal/search/algolia"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/pkg/types"
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

	if err = im.Wipe(ctx); err != nil {
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
	filter.Limit = pointers.Pointer(uint8(50))

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
